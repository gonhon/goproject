package rpc

import (
	"encoding/json"
	"log"
	"net"
	"reflect"
	"sync"

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

type Server struct{}

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
	f, ok := codec.NewCodeFuncMap[opt.CodecType]
	if ok {
		log.Printf("rpc server invalid CodecType %s", opt.CodecType)
		return
	}
	server.serverCodec(f(conn))

}

var invalidRequest = struct{}{}

func (server *Server) serverCodec(c codec.Codec) {
	//TODO
	/* sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {

	} */

}

//---------------------serverCodec:读取、处理、回复---------------------
type request struct {
	h            codec.Header
	argv, replyv reflect.Value
}

//读取数据
func (server *Server) readRequest(c codec.Codec) (*request, error) {
	return nil, nil
}

//获取head信息
func (server *Server) readRequestHeader(c codec.Codec) (*codec.Header, error) {
	return nil, nil
}

//数据响应
func (server *Server) sendResponse(c codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {

}

//数据处理
func (server *Server) handleRequest(c codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {

}
