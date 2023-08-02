package main

import (
	"fmt"
	"net"
)

func main() {
	//使用 net.Listen 监听连接的地址与端口
	listener, err := net.Listen("tcp", "0.0.0.0:9091")
	if err != nil {
		fmt.Printf("listen fail, err: %v\n", err)
		return
	}
	for {
		// 等待连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept fail, err: %v\n", err)
			continue
		}
		// 对每个新连接创建一个协程进行收发数据
		go process(conn)
	}
}
func process(conn net.Conn) {
	defer conn.Close()
	for {
		var buf [128]byte
		//接受数据
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("read from connect failed, err: %v\n", err)
			break
		}
		fmt.Printf("receive from client, data: %v\n", string(buf[:n]))
		fmt.Printf("receive from client, data: %v\n", buf[:n])
		//发送数据
		if _, err = conn.Write([]byte("Send From Server")); err != nil {
			fmt.Printf("write to client failed, err: %v\n", err)
			break
		}
	}
}
