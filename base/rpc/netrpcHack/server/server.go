package main

import (
	"fmt"
	"github.com/limerence-code/goproject/base/rpc/rpcIntor"
	"log"
	"math/rand"
	"net"
	"net/rpc"
	"sync"
	"time"
)

type KVStoreService struct {
	m      map[string]string
	filter map[string]func(string)
	mux    sync.Mutex
}

func NewKVStoreService() *KVStoreService {
	return &KVStoreService{
		m:      make(map[string]string),
		filter: make(map[string]func(string)),
	}
}

func (p *KVStoreService) Get(key string, value *string) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if v, ok := p.m[key]; ok {
		*value = v
		return nil
	}
	return fmt.Errorf("not found")
}

func (p *KVStoreService) Set(kv [2]string, reply *struct{}) error {
	p.mux.Lock()
	defer p.mux.Unlock()
	key, value := kv[0], kv[1]
	if oldVal := p.m[key]; oldVal != value {
		for _, filter := range p.filter {
			filter(key)
		}
	}
	p.m[key] = value
	log.Println("set:key,val", key, value)
	return nil
}

func (p *KVStoreService) Watch(timeoutSecond int, keyChanged *string) error {
	id := fmt.Sprintf("watch-%s-%03d", time.Now(), rand.Int())
	chans := make(chan string, 10)
	p.mux.Lock()
	p.filter[id] = func(key string) {
		chans <- key
	}
	defer p.mux.Unlock()

	select {
	case <-time.After(time.Duration(timeoutSecond) * time.Second):
		log.Println("超时...")
		return fmt.Errorf("timeout")
	case key := <-chans:
		*keyChanged = key
		log.Println("获取到key：", key)
	}

	return nil
}

func main() {
	err := rpc.RegisterName(rpcIntor.WatchServiceName, NewKVStoreService())
	if err != nil {
		return
	}
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
		go rpc.ServeConn(conn)
		//go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
