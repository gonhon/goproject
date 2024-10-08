package main

import (
	"fmt"
	"net"

	"github.com/gonhon/tcp-server-demo1/frame"
	"github.com/gonhon/tcp-server-demo1/packet"
)

func handeConn(c net.Conn) {
	defer c.Close()

	frameCodec := frame.NewMyFrameCodec()
	for {
		framePayload, err := frameCodec.Decode(c)
		if err != nil {
			fmt.Println("handleConn:frame decode error:", err)
			return
		}
		ackFramePayload, err := handlePacket(framePayload)
		if err != nil {
			fmt.Println("hanldeConn:handle packet error:", err)
			return
		}
		err = frameCodec.Encode(c, ackFramePayload)
		if err != nil {
			fmt.Println("handleConn:frame encode error:", err)
			return
		}
	}
}

func handlePacket(framePayload []byte) (ackFramePayload []byte, err error) {
	var p packet.Packet
	p, err = packet.Decode(framePayload)
	if err != nil {
		fmt.Println("handleConn:packet decode error:", err)
		return
	}
	switch p.(type) {
	case *packet.Submit:
		submit := p.(*packet.Submit)
		fmt.Printf("revc submit:id=%s,payload=%s\n", submit.ID, string(submit.Payload))
		submitAck := &packet.SubmitAck{ID: submit.ID, Result: 0}
		ackFramePayload, err = packet.Encode(submitAck)
		if err != nil {
			fmt.Println("handleConn: packet encode error:", err)
			return nil, err
		}
		return ackFramePayload, nil
	default:
		return nil, fmt.Errorf("unknow packet type")
	}
}

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	fmt.Println("server start ok (on *.8888)")
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accpect error:", err)
			break
		}
		go handeConn(c)
	}

}
