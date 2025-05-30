// p2p-client/main.go
package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type Message struct {
	Type      string          `json:"type"`
	Room      string          `json:"room"`
	SDP       json.RawMessage `json:"sdp,omitempty"`
	Candidate json.RawMessage `json:"candidate,omitempty"`
}

var (
	peerConnection *webrtc.PeerConnection
	dataChannel    *webrtc.DataChannel
	wsConn         *websocket.Conn
	roomID         string
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// 输入房间号
	log.Print("Enter room ID: ")
	room, _ := reader.ReadString('\n')
	roomID = strings.TrimSpace(room)

	// 连接信令服务器
	connectSignalingServer()

	// 配置WebRTC
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
			// 如果需要TURN服务器，取消注释并填写你的TURN服务器信息
			// {
			// 	URLs:       []string{"turn:your-turn-server.example"},
			// 	Username:   "username",
			// 	Credential: "password",
			// },
		},
	}

	var err error
	peerConnection, err = webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatal("Failed to create peer connection:", err)
	}

	// 设置ICE候选处理
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		candidateJSON, err := json.Marshal(c.ToJSON())
		if err != nil {
			log.Println("Marshal ICE candidate error:", err)
			return
		}

		msg := Message{
			Type:      "candidate",
			Room:      roomID,
			Candidate: candidateJSON,
		}

		if err := wsConn.WriteJSON(msg); err != nil {
			log.Println("Send ICE candidate error:", err)
		}
	})

	// 创建数据通道
	dataChannel, err = peerConnection.CreateDataChannel("chat", nil)
	if err != nil {
		log.Fatal("Failed to create data channel:", err)
	}

	dataChannel.OnOpen(func() {
		log.Println("Data channel opened!")
	})

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		log.Printf("Received: %s\n", string(msg.Data))
	})

	// 处理远程数据通道
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		log.Println("New DataChannel:", d.Label())

		dataChannel = d
		dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
			log.Printf("Received: %s\n", string(msg.Data))
		})
	})

	// 创建offer
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		log.Fatal("Failed to create offer:", err)
	}

	if err = peerConnection.SetLocalDescription(offer); err != nil {
		log.Fatal("Failed to set local description:", err)
	}

	// 发送offer
	offerJSON, err := json.Marshal(offer)
	if err != nil {
		log.Fatal("Marshal offer error:", err)
	}

	msg := Message{
		Type: "offer",
		Room: roomID,
		SDP:  offerJSON,
	}

	if err := wsConn.WriteJSON(msg); err != nil {
		log.Fatal("Send offer error:", err)
	}

	// 启动消息输入循环
	go readMessages()

	// 保持程序运行
	select {}
}

func connectSignalingServer() {
	var err error
	wsConn, _, err = websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}

	// 加入房间
	msg := Message{
		Type: "join",
		Room: roomID,
	}

	if err := wsConn.WriteJSON(msg); err != nil {
		log.Fatal("Join room error:", err)
	}

	// 处理信令消息
	go func() {
		for {
			var msg Message
			if err := wsConn.ReadJSON(&msg); err != nil {
				log.Println("Read message error:", err)
				return
			}

			switch msg.Type {
			case "offer":
				handleOffer(msg.SDP)
			case "answer":
				handleAnswer(msg.SDP)
			case "candidate":
				handleCandidate(msg.Candidate)
			}
		}
	}()
}

func handleOffer(sdp json.RawMessage) {
	offer := webrtc.SessionDescription{}
	if err := json.Unmarshal(sdp, &offer); err != nil {
		log.Println("Unmarshal offer error:", err)
		return
	}

	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		log.Println("Set remote description error:", err)
		return
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		log.Println("Create answer error:", err)
		return
	}

	if err = peerConnection.SetLocalDescription(answer); err != nil {
		log.Println("Set local description error:", err)
		return
	}

	answerJSON, err := json.Marshal(answer)
	if err != nil {
		log.Println("Marshal answer error:", err)
		return
	}

	msg := Message{
		Type: "answer",
		Room: roomID,
		SDP:  answerJSON,
	}

	if err := wsConn.WriteJSON(msg); err != nil {
		log.Println("Send answer error:", err)
	}
}

func handleAnswer(sdp json.RawMessage) {
	answer := webrtc.SessionDescription{}
	if err := json.Unmarshal(sdp, &answer); err != nil {
		log.Println("Unmarshal answer error:", err)
		return
	}

	if err := peerConnection.SetRemoteDescription(answer); err != nil {
		log.Println("Set remote description error:", err)
	}
}

func handleCandidate(candidate json.RawMessage) {
	iceCandidate := webrtc.ICECandidateInit{}
	if err := json.Unmarshal(candidate, &iceCandidate); err != nil {
		log.Println("Unmarshal ICE candidate error:", err)
		return
	}

	if err := peerConnection.AddICECandidate(iceCandidate); err != nil {
		log.Println("Add ICE candidate error:", err)
	}
}

func readMessages() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if dataChannel != nil && dataChannel.ReadyState() == webrtc.DataChannelStateOpen {
			if err := dataChannel.SendText(text); err != nil {
				log.Println("Send message error:", err)
			}
		} else {
			log.Println("Data channel not ready")
		}
	}
}
