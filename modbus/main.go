package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/grid-x/modbus"
)

func main() {
	fmt.Println(time.Now().UnixMilli())
	fmt.Println(time.Now().Unix())
	//Modbus()

	//var client *modbus2.ModbusClient
	//var err error

	// for a TCP endpoint
	// (see examples/tls_client.go for TLS usage and options)
	//client, err = modbus2.NewClient(&modbus2.ClientConfiguration{
	//	URL:     "tcp://192.168.1.88:502", // 可以直接使用ngrok映射到外网的地址
	//	Timeout: 1 * time.Second,
	//})
	// note: use udp:// for modbus TCP over UDP
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//err = client.Open()
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//address := uint16(1)
	//v, err := client.ReadRegister(address, modbus2.HOLDING_REGISTER)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println("address:", address, ",old value:", v)

	// uint
	//fmt.Println("uint read and write")
	//for i := 0; i <= 7; i++ {
	//	address := uint16(i)
	//	value := uint16(1)
	//	v, err := client.ReadRegister(address, modbus2.HOLDING_REGISTER)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//	fmt.Println("address:", address, ",old value:", v)
	//	err = client.WriteRegister(address, value)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//	v, err = client.ReadRegister(address, modbus2.HOLDING_REGISTER)
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}
	//	fmt.Println("address:", address, ",new value:", v)
	//}

	//// float
	//fmt.Println("float read and write")
	//bytes, err := client.ReadBytes(6, 4, modbus2.HOLDING_REGISTER)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//f, _ := bytesToFloat32(bytes)
	//fmt.Println("old value:", f)
	//err = client.WriteBytes(6, float32ToBytes(123.456))
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//bytes, err = client.ReadBytes(6, 4, modbus2.HOLDING_REGISTER)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//f, _ = bytesToFloat32(bytes)
	//fmt.Println("new value:", f)

}

func Modbus() {
	var ModeBusTCpClient *modbus.TCPClientHandler
	ip := "192.168.1.88:502"

	ModeBusTCpClient = modbus.NewTCPClientHandler(ip)
	err := ModeBusTCpClient.Connect()
	fmt.Println(err)
	defer ModeBusTCpClient.Close()
	client := modbus.NewClient(ModeBusTCpClient)

	address := uint16(3)                                        // D3 40004 // 最前面的4表明它是一个保持寄存器，后面的0004指第四个保持存储器，并且它的地址是3
	coil, err := client.WriteSingleRegister(address, uint16(2)) // 写单个保持寄存器
	fmt.Println("ip: ", ip)
	fmt.Println("address: ", address)

	if err != nil {
		fmt.Printf("err is %v", err.Error())
		return
	}
	fmt.Println(int(binary.BigEndian.Uint16(coil)) == 2)
}
