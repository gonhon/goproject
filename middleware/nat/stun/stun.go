package main

import (
	"fmt"
	"log"

	"github.com/pion/stun"
)

func main() {
	// 配置 STUN 服务器（如 Google 公共服务器）
	stunServer := "stun.l.google.com:19302"

	// 创建 STUN 客户端
	c, err := stun.Dial("udp", stunServer)
	if err != nil {
		log.Fatal("STUN 连接失败:", err)
	}
	defer c.Close()

	// 发送 Binding Request 获取公网地址
	var publicAddr stun.XORMappedAddress
	if err := c.Do(stun.MustBuild(stun.TransactionID, stun.BindingRequest), func(res stun.Event) {
		if res.Error != nil {
			log.Fatal("STUN 请求失败:", res.Error)
		}
		// 解析返回的公网地址
		if err := publicAddr.GetFrom(res.Message); err != nil {
			log.Fatal("解析地址失败:", err)
		}
	}); err != nil {
		log.Fatal("STUN 操作失败:", err)
	}

	fmt.Printf("公网地址: %s\n", publicAddr)
}
