package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	topic := "transmission_line_camera"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "172.20.50.120:29092", topic, partition)
	if err != nil {
		fmt.Println("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		fmt.Println("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		fmt.Println("failed to close writer:", err)
	}
}
