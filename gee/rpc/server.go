package rpc

import (
	"encoding/json"
	"errors"
	"go/ast"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/limerence-code/goproject/gee/rpc/codec"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int
	CodecType   codec.Type
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
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
	server.serverCodec(f(conn))

}

var invalidRequest = struct{}{}

func (server *Server) serverCodec(c codec.Codec) {
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
		go server.handleRequest(c, req, sending, wg)
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
func (server *Server) handleRequest(c codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	err := req.svc.call(req.mtype, req.argv, req.replyv)
	if err != nil {
		req.h.Error = err.Error()
		server.sendResponse(c, req.h, invalidRequest, sending)
		return
	}

	// log.Println(req.h, req.argv.Elem())

	server.sendResponse(c, req.h, req.replyv.Interface(), sending)
}

//注册服务到map
func (server *Server) Regisert(rcvr interface{}) error {
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
	return DefaultServer.Regisert(rcvr)
}

type methedType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	numCalls  uint64
}

func (m *methedType) NumCalls() uint64 {
	return atomic.LoadUint64(&m.numCalls)
}

func (m *methedType) newArgv() reflect.Value {
	var argv reflect.Value
	//指针类型
	if m.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(m.ArgType.Elem())
	} else {
		//值类型
		argv = reflect.New(m.ArgType).Elem()
	}
	return argv
}
func (m *methedType) newReplyv() reflect.Value {
	replyv := reflect.New(m.ReplyType.Elem())
	switch m.ReplyType.Elem().Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}
	return replyv
}

type service struct {
	//映射的结构体的名称
	name string
	//结构体的类型
	typ reflect.Type
	//结构体的实例本身
	rcvr reflect.Value
	//结构体的所有符合条件的方法
	method map[string]*methedType
}

func newService(rcvr interface{}) *service {
	s := new(service)

	s.rcvr = reflect.ValueOf(rcvr)
	s.name = reflect.Indirect(s.rcvr).Type().Name()
	s.typ = reflect.TypeOf(rcvr)
	if !ast.IsExported(s.name) {
		log.Fatalf("rpc server: %s is not a valid service name", s.name)
	}
	s.registerMethods()

	return s

}
func (s *service) registerMethods() {
	s.method = make(map[string]*methedType)
	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i)
		mType := method.Type
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		argType, replyType := mType.In(1), mType.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}
		s.method[method.Name] = &methedType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
		log.Printf("rpc server: register %s.%s\n", s.name, method.Name)
	}

}

func (s *service) call(m *methedType, argv, replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls, 1)
	f := m.method.Func
	returnValues := f.Call([]reflect.Value{s.rcvr, argv, replyv})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil

}

func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}
