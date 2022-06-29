package main

import (
	"fmt"
	"github.com/limerence-code/goproject/base/rpc/protos"
	"github.com/limerence-code/goproject/base/rpc/rpcIntor"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloServiceClient struct {
	*rpc.Client
}

// var  HelloServiceInterface = (*HelloServiceClient)(nil)

func DailHelloService(network, address string) (*HelloServiceClient, error) {
	// client, err := rpc.Dial(network, address)

	//使用net获取
	conn, err := net.Dial(network, address)
	//进行包装
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	if err != nil {
		log.Fatal("rpc err:", err)
		return nil, err
	}
	return &HelloServiceClient{client}, nil
}

func (h *HelloServiceClient) Hello(request *protos.String, reply *protos.String) error {
	return h.Call(rpcIntor.HelloServiceName+".Hello", request, reply)
}

func main() {
	reply := &protos.String{Value: " go "}

	helloServer, err := DailHelloService("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("rpc err:", err)
	}
	helloServer.Hello(&protos.String{Value: "go"}, reply)

	// client, err := rpc.Dial("tcp", "127.0.0.1:1234")
	// if err != nil {
	// 	log.Fatal("rpc err:", err)
	// }

	// //调用远程接口
	// client.Call(rpcIntor.HelloServiceName+".Hello", "hello", &reply)
	fmt.Println("reply===>:", reply.GetValue())
}
