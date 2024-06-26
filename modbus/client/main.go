package main

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/grid-x/modbus"
)

func main() {
	var ModeBusTCpClient *modbus.TCPClientHandler
	//ip := "127.0.0.1:1502"
	ip := "192.168.1.2:502"
	ModeBusTCpClient = modbus.NewTCPClientHandler(ip)
	ModeBusTCpClient.LinkRecoveryTimeout = 1 * time.Second
	ModeBusTCpClient.ProtocolRecoveryTimeout = 1 * time.Second
	// 设置从机地址
	//ModeBusTCpClient.SetSlave(byte(1))
	// 连接 Modbus server (slave)
	err := ModeBusTCpClient.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ModeBusTCpClient.Close()
	client := modbus.NewClient(ModeBusTCpClient)
	// WriteSingleRegister 写单路寄存器 对应功能码06
	// address 寄存器地址 // value 要写入的值

	for {
		coil, err := client.ReadHoldingRegisters(2, uint16(1))
		if err != nil {
			fmt.Printf("err is %s\n", err.Error())
			continue
		}
		fmt.Println(int(binary.BigEndian.Uint16(coil)))
		time.Sleep(1 * time.Second)
	}
}
