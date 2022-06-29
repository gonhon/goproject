package tcp

import (
	"fmt"
	"net"
)

func TcpServer() {
	fmt.Println("服务端启动开始监听端口")
	listen, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer listen.Close()

	for {
		conn, e := listen.Accept()
		if e != nil {
			fmt.Println("Accept", e)
		} else {

		}
		go process(conn)
	}
}
func process(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		fmt.Println("等待客户端来连接：", conn.RemoteAddr().String())
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read err-->", err)
			return
		}
		fmt.Println(string(buf[:n]))
	}

}
