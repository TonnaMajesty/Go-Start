package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有follower都回复ack，确保Kafka不会丢消息
	config.Producer.Return.Successes = true
	//config.Producer.Partitioner = sarama.NewHashPartitioner // 对Key进行Hash，同样的Key每次都落到一个分区，这样消息是有序的

	// 使用同步producer，异步模式下有更高的性能，但是处理更复杂，这里建议先从简单的入手
	producer, err := sarama.NewSyncProducer([]string{"172.20.50.120:29092"}, config)
	defer func() {
		_ = producer.Close()
	}()

	if err != nil {
		panic(err.Error())
	}

	// 模拟4个消息
	//for {

	//data := DeviceData{
	//	Devices: []Device{
	//		{
	//			DeviceID: "device1",
	//			Services: []Service{
	//				{
	//					ServiceID: "1",
	//					Data: map[string]interface{}{
	//						"SN":                 "SN1",
	//						"Name":               "Name1",
	//						"OperationTeam":      "OperationTeam1",
	//						"Line":               "Line1",
	//						"Tower":              "Tower1",
	//						"DischargeVoltage":   1,
	//						"WorkTemperature":    1,
	//						"BatteryQuantity":    1,
	//						"BatteryState":       1,
	//						"ContinuousWorkTime": 1,
	//						"InstallLocation":    1,
	//					},
	//					EventTime: "",
	//				},
	//			},
	//		},
	//	},
	//}

	pic, err := ioutil.ReadFile("/Users/tonnamajesty/Downloads/bottle_front.jpg")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	picBase64str := base64.StdEncoding.EncodeToString(pic)

	data := DeviceData{
		Devices: []Device{
			{
				DeviceID: "device1",
				Services: []Service{
					{
						ServiceID: "2",
						Data: map[string]interface{}{
							"timeStamp": time.Now().UnixMilli(), // 时间戳
							"height":    2448,                   // 高度
							"width":     2048,                   // 宽度
							"data":      picBase64str,           // base64
						},
						EventTime: "20191031T144831Z",
					},
				},
			},
		},
	}

	dataStr, _ := json.Marshal(data)

	//rand.Seed(int64(time.Now().Nanosecond()))
	msg := &sarama.ProducerMessage{
		Topic: "transmission_line_camera",
		Value: sarama.StringEncoder(dataStr),
		Key:   sarama.StringEncoder("BBB"),
	}

	t1 := time.Now().Nanosecond()
	partition, offset, err := producer.SendMessage(msg)
	t2 := time.Now().Nanosecond()

	if err == nil {
		fmt.Println("produce success, partition:", partition, ",offset:", offset, ",cost:", (t2-t1)/(1000*1000), " ms")
	} else {
		fmt.Println(err.Error())
	}

	//time.Sleep(1 * time.Second)
	//}
}

const (
	topic   = "transmission_line_camera"
	groupID = "transmission_line_consumer_group"
)

type DeviceData struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	DeviceID string    `json:"deviceId"`
	Services []Service `json:"services"`
}

type Service struct {
	Data      map[string]interface{} `json:"data"`
	EventTime string                 `json:"eventTime,omitempty"`
	ServiceID string                 `json:"serviceId"`
}

type DeviceInfo struct {
	SN                 string // 摄像头序列号
	Name               string // 摄像头名称
	OperationTeam      string // 运维班组
	Line               string // 输电线路
	Tower              string // 杆塔编号
	DischargeVoltage   int    // 放电电压
	WorkTemperature    int    // 工作温度
	BatteryQuantity    int    // 电池电量
	BatteryState       int    // 电池放电状态
	ContinuousWorkTime int    // 连续工作时间
	InstallLocation    int    // 安装位置
}

type DevicePhoto struct {
	TimeStamp int    // 时间戳
	Height    int    // 高度
	Weight    int    // 宽度
	Data      string // base64
}
