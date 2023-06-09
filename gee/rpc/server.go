package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/limerence-code/goproject/gee/rpc/codec"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber    int
	CodecType      codec.Type
	ConnectTimeout time.Duration
	HandleTimeout  time.Duration
}

var DefaultOption = &Option{
	MagicNumber:    MagicNumber,
	CodecType:      codec.GobType,
	ConnectTimeout: time.Second * 10,
}

type Server struct {
	ServiceMap sync.Map
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (server *Server) Accept(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("rpc server error: ", err)
		}
		go server.ServerConn(conn)
	}
}

func Accept(listener net.Listener) {
	DefaultServer.Accept(listener)
	// ls,_:=net.Listen("tcp", ":80")
	// Accept(ls)
}

func (server *Server) ServerConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server option error:", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server invalid MagicNumber %x", opt.MagicNumber)
		return
	}
	f, ok := codec.NewCodecFuncMap[opt.CodecType]
	if !ok {
		log.Printf("rpc server invalid CodecType %s", opt.CodecType)
		return
	}
	server.serverCodec(f(conn), &opt)

}

var invalidRequest = struct{}{}

func (server *Server) serverCodec(c codec.Codec, opt *Option) {
	//TODO
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := server.readRequest(c)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			server.sendResponse(c, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(c, req, sending, wg, opt.HandleTimeout)
	}
	wg.Wait()
	c.Close()

}

// ---------------------serverCodec:读取、处理、回复---------------------
type request struct {
	h            *codec.Header
	argv, replyv reflect.Value
	mtype        *methedType
	svc          *service
}

// 读取数据
func (server *Server) readRequest(c codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(c)
	if err != nil {
		return nil, err
	}
	req := &request{
		h: h,
	}

	req.svc, req.mtype, err = server.findService(h.ServiceMethod)
	if err != nil {
		return req, err
	}

	req.argv = req.mtype.newArgv()
	req.replyv = req.mtype.newReplyv()

	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}

	if err = c.ReadBody(argvi); err != nil {
		log.Println("rpc server: read argv err:", err)
		return req, nil
	}
	return req, nil
}

// 获取head信息
func (server *Server) readRequestHeader(c codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := c.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Panicln("rpc server read header error", err)
		}
		return nil, err
	}
	return &h, nil
}

// 数据响应
func (server *Server) sendResponse(c codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()

	if err := c.Write(h, body); err != nil {
		log.Println("rpc server  write response error:", err)
	}
}

// 数据处理
func (server *Server) handleRequest(c codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()

	callChan := make(chan struct{})
	sendChan := make(chan struct{})
	go func() {
		err := req.svc.call(req.mtype, req.argv, req.replyv)
		callChan <- struct{}{}
		//异常
		if err != nil {
			req.h.Error = err.Error()
			server.sendResponse(c, req.h, invalidRequest, sending)
			sendChan <- struct{}{}
			return
		}
		server.sendResponse(c, req.h, req.replyv.Interface(), sending)
		sendChan <- struct{}{}
	}()

	if timeout == 0 {
		<-callChan
		<-sendChan
		return
	}
	select {
	case <-time.After(timeout):
		req.h.Error = fmt.Sprintf("rpc server: request handle timeout: expect within %s", timeout)
		server.sendResponse(c, req.h, invalidRequest, sending)
	case <-callChan:
		<-sendChan
	}
}

// 注册服务到map
func (server *Server) Register(rcvr interface{}) error {
	s := newService(rcvr)
	if _, dup := server.ServiceMap.LoadOrStore(s.name, s); dup {
		return errors.New("rpc:service already defined:" + s.name)
	}
	return nil
}
func (server *Server) findService(serviceMethod string) (svc *service, mtype *methedType, err error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc server: service/method request ill-formed: " + serviceMethod)
	}
	//获取服务名和方法名
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	svci, ok := server.ServiceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc server: can't find service " + serviceName)
		return
	}
	svc = svci.(*service)
	mtype = svc.method[methodName]
	if mtype == nil {
		err = errors.New("rpc server: can't find method " + methodName)
	}
	return
}

func Register(rcvr interface{}) error {
	return DefaultServer.Register(rcvr)
}

//===============================http===============================

const (
	connected        = "200 connected to rpc"
	defaultRpcPath   = "/rpc"
	defaultDebugPath = "/debug/rpc"
)

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "CONNECT" {
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "405 must CONNECT\n")
		return
	}
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("rpc hijack error", r.RemoteAddr, ":", err.Error())
		return
	}
	io.WriteString(conn, "HTTP/1.0 "+connected+"\n\n")
	server.ServerConn(conn)
}

func (server *Server) HandleHTTP() {
	http.Handle(defaultRpcPath, server)
	http.Handle(defaultDebugPath, debugHTTP{server})
	log.Println("rpc server debug path:", defaultDebugPath)
}

func HandleHTTP() {
	DefaultServer.HandleHTTP()
}
