package main

import (
	"fmt"
	"github.com/limerence-code/goproject/base/rpc/rpcIntor"
	"log"
	"net/rpc"
	"time"
)

func doClientWork(client *rpc.Client) {
	go func() {
		var keyChanged string
		err := client.Call(rpcIntor.WatchServiceName+".Watch", 30, &keyChanged)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("watch:", keyChanged)
	}()

	time.Sleep(time.Second * 1)

	err := client.Call(rpcIntor.WatchServiceName+".Set", [2]string{"abc", "abc-value"}, new(struct{}))
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 3)
}
func main() {
	//conn, err := net.Dial("tcp", "127.0.0.1:1234")

	client, err := rpc.Dial("tcp", "127.0.0.1:1234")

	//进行包装
	//client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	if err != nil {
		log.Fatal("rpc err:", err)
	} else {
		doClientWork(client)
	}
}
