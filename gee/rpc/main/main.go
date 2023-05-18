package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/limerence-code/goproject/gee/rpc"
	"github.com/limerence-code/goproject/gee/rpc/codec"
)

func startServer(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error: ", err)
	}
	log.Println("start rpc server on ", l.Addr())
	addr <- l.Addr().String()
	rpc.Accept(l)
}

func main() {
	addr := make(chan string)
	go startServer(addr)

	conn, _ := net.Dial("tcp", <-addr)
	defer func() { conn.Close() }()
	time.Sleep(time.Second)
	json.NewEncoder(conn).Encode(rpc.DefaultOption)
	c := codec.NewGobCodec(conn)

	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		c.Write(h, fmt.Sprintf("rpc request %d", h.Seq))
		c.ReadHeader(h)
		var reply string
		c.ReadBody(&reply)
		log.Println("reply: ", reply)
	}
}
