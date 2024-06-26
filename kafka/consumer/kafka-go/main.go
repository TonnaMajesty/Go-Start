package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func main() {
	// Kafka broker连接信息
	brokerUrls := []string{"172.20.50.120:29092"}
	//brokerUrls := []string{"localhost:9092"}
	topic := "transmission_line_camera"
	//groupID := "my-group"

	// Kafka consumer配置
	config := kafka.ReaderConfig{
		Brokers:         brokerUrls,
		Topic:           topic,
		GroupID:         "consumer-group-id",
		MinBytes:        1,                // 最小读取字节数
		MaxBytes:        10e6,             // 最大读取字节数
		MaxWait:         10 * time.Second, // 最大等待时间
		ReadLagInterval: -1,               // 不计算消费滞后
	}

	// 创建Kafka reader
	r := kafka.NewReader(config)

	// 捕获系统中断信号
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// 创建context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动消费者
	go func() {
		for {
			m, err := r.FetchMessage(ctx)
			if err != nil {
				break
			}
			fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
			if err := r.CommitMessages(ctx, m); err != nil {
				log.Fatal("failed to commit messages:", err)
			}
		}
	}()

	// 等待系统中断信号
	<-signals

	// 关闭reader和context
	r.Close()
	cancel()
}
