package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 将 HTTP 连接升级为 WebSocket 连接
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade failed:", err)
			return
		}

		// 启动探活 Goroutine
		stopCh := make(chan struct{})
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					// 向客户端发送 Ping 消息
					err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
					if err != nil {
						if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
							log.Println("Connection closed normally")
						} else {
							log.Println("WriteControl failed:", err)
						}
						close(stopCh)
						return
					}
				case <-stopCh:
					return
				}
			}
		}()

		go func() {
			// 关闭 WebSocket 连接
			defer func() {
				err = conn.Close()
				if err != nil {
					log.Println("Close failed:", err)
					return
				}
			}()

			// 循环读取 WebSocket 消息
			for {
				select {
				case <-stopCh:
					return
				default:
					messageType, message, err := conn.ReadMessage()
					if err != nil {
						if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
							log.Println("Connection closed normally")
						} else {
							log.Println("ReadMessage failed:", err)
						}
						return
					}
					log.Println("Received message:", string(message))

					// 向客户端发送响应消息
					err = conn.WriteMessage(messageType, message)
					if err != nil {
						log.Println("WriteMessage failed:", err)
						return
					}
				}
			}
		}()

	})

	// 启动 HTTP 服务器
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
