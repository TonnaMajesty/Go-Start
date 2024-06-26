package main

import (
	"log"

	"github.com/Shopify/sarama"
)

func main() {
	// 设置 Kafka Consumer 的配置信息
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	// 创建 Kafka Consumer 连接
	client, err := sarama.NewClient([]string{"localhost:9093"}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	// 指定要消费的 Topic 和 Partition
	topic := "test"
	partition := int32(0)

	// 获取 Partition 的 Offset
	offsetManager, err := sarama.NewOffsetManagerFromClient("my-group", client)
	if err != nil {
		log.Fatal(err)
	}
	defer offsetManager.Close()

	partitionOffsetManager, err := offsetManager.ManagePartition(topic, partition)
	if err != nil {
		log.Fatal(err)
	}
	defer partitionOffsetManager.Close()

	offset, _ := partitionOffsetManager.NextOffset()

	// 创建 Partition Consumer
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer partitionConsumer.Close()

	// 循环读取消息
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Message received: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s\n",
				msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			partitionOffsetManager.MarkOffset(msg.Offset+1, "")
		case err := <-partitionConsumer.Errors():
			log.Printf("Error received: %v\n", err)
		}
	}
}
