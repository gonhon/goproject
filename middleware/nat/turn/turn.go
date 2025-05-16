package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/pion/turn/v2"
)

func main() {
	// TURN 服务器配置
	turnServer := "your-turn-server.com:3478"
	username := "your-username"
	password := "your-password"

	// 创建 UDP 连接
	conn, err := net.ListenPacket("udp", "0.0.0.0:0")
	if err != nil {
		log.Fatal("UDP 监听失败:", err)
	}
	defer conn.Close()

	// 配置 TURN 客户端
	turnClient, err := turn.NewClient(&turn.ClientConfig{
		STUNServerAddr: turnServer,
		TURNServerAddr: turnServer,
		Conn:           conn,
		Username:       username,
		Password:       password,
		Realm:          "your-realm",
	})
	if err != nil {
		log.Fatal("TURN 客户端创建失败:", err)
	}
	defer turnClient.Close()

	// 分配中继地址
	relayAddr, err := turnClient.Allocate()
	if err != nil {
		log.Fatal("TURN 分配中继地址失败:", err)
	}
	fmt.Printf("中继地址: %s\n", relayAddr)

	// 保持连接（防止 NAT 超时）
	for {
		time.Sleep(5 * time.Second)
		if _, err := turnClient.SendBindingRequest(); err != nil {
			log.Fatal("保活失败:", err)
		}
	}
}
