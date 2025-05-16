package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pion/ice/v2"
)

func main() {
	// 1. 创建 ICE Agent
	agent, err := ice.NewAgent(&ice.AgentConfig{
		NetworkTypes: []ice.NetworkType{ice.NetworkTypeUDP4},
		Urls: []*ice.URL{
			{ // 公共 STUN 服务器
				Scheme: ice.SchemeTypeSTUN,
				Host:   "stun.l.google.com",
				Port:   19302,
			},
			{ // 自建 TURN 服务器（可选）
				Scheme:   ice.SchemeTypeTURN,
				Host:     "your-turn.example.com",
				Port:     3478,
				Username: "user",
				Password: "pass",
			},
		},
	})
	if err != nil {
		log.Fatal("Failed to create ICE agent:", err)
	}

	// 2. 监听 ICE 连接状态变化
	agent.OnConnectionStateChange(func(state ice.ConnectionState) {
		fmt.Printf("ICE Connection State: %s\n", state)
		if state == ice.ConnectionStateConnected {
			fmt.Println("ICE 连接成功！")
		} else if state == ice.ConnectionStateFailed {
			log.Fatal("ICE 连接失败")
		}
	})

	// 3. 收集本地候选地址（Local Candidates）
	if err := agent.GatherCandidates(); err != nil {
		log.Fatal("Failed to gather candidates:", err)
	}

	// 4. 获取本地候选地址（需通过信令服务器发送给对端）
	localCandidates, err := agent.GetLocalCandidates()
	if err != nil {
		log.Fatal("Failed to get local candidates:", err)
	}
	for _, c := range localCandidates {
		fmt.Printf("Local Candidate: %s\n", c)
	}

	// 5. 模拟接收对端候选地址（实际需通过信令服务器交换）
	// remoteCandidates := []ice.Candidate{} // 替换为对端的 Candidates
	var remoteCandidate ice.Candidate // 替换为对端的 Candidates
	if err := agent.AddRemoteCandidate(remoteCandidate); err != nil {
		log.Fatal("Failed to set remote candidates:", err)
	}

	// 6. 设置超时（例如 10 秒）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 7. 等待连接完成或超时
	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			log.Fatal("ICE 连接超时")
		}
	}
}
