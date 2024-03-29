package main

import (
	"fmt"
	"github.com/limerence-code/goproject/base/rpc/rpcIntor"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloServiceInterface interface {
	Hello(req string, reply *string) error
}

type HelloService struct {
}

func RegisterHelloService(server HelloServiceInterface) {
	rpc.RegisterName(rpcIntor.HelloServiceName, server)
}

func (helloService *HelloService) Hello(req string, reply *string) error {
	*reply = req + " go"
	return nil
}

func main() {
	rpcCon()
	// httpCon()
}

func rpcCon() {
	//rpc注册服务
	// rpc.RegisterName("HelloService", new(HelloService))
	RegisterHelloService(new(HelloService))
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTcp error:", err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal("AcceptTcp  error:", err)
		}
		fmt.Println("rpc conn...")
		// go rpc.ServeConn(conn)
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}

}

// 支持http
//
//	$ curl localhost:1234/jsonrpc -X POST \
//		--data '{"method":"HelloService.Hello","params":["hello"],"id":0}'
func httpCon() {
	rpc.RegisterName("HelloService", new(HelloService))

	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}
		fmt.Println("http conn...")
		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})

	http.ListenAndServe(":1234", nil)
}
