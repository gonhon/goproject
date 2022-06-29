package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func tcpClient() {
	fmt.Println("客户端")
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Println("conn-->", err)
		return
	}

	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read-->", err)
		}

		if line = strings.Trim(line, "\r\n"); line == "exit" {
			fmt.Println("客户端准备退出")
			return
		}

		_, err = conn.Write([]byte(line + "\n"))
		if err != nil {
			fmt.Println("write-->", err)
		}
	}
}

func main() {
	tcpClient()
}
