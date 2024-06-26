package main

import (
	"encoding/binary"
	"fmt"

	"testing"
	"time"

	"github.com/grid-x/modbus"
	"github.com/tbrandon/mbserver"
)

func TestModbus(t *testing.T) {
	var ModeBusTCpClient *modbus.TCPClientHandler
	//ip := "127.0.0.1:1502"
	ip := "192.168.1.2:502"
	ModeBusTCpClient = modbus.NewTCPClientHandler(ip)
	// 设置从机地址
	ModeBusTCpClient.SetSlave(byte(1))
	// 连接 Modbus server (slave)
	ModeBusTCpClient.Connect()
	defer ModeBusTCpClient.Close()
	client := modbus.NewClient(ModeBusTCpClient)
	// WriteSingleRegister 写单路寄存器 对应功能码06
	// address 寄存器地址 // value 要写入的值

	coil, err := client.WriteSingleRegister(1, uint16(1))
	fmt.Println(err)
	if err != nil {
		fmt.Errorf("err is %v", err.Error())
		return
	}
	fmt.Println(int(binary.BigEndian.Uint16(coil)) == 2)
}

// creates a new Modbus server (slave).
func TestSimulateModbusTcpServer(t *testing.T) {
	serv := mbserver.NewServer()
	err := serv.ListenTCP("127.0.0.1:1502")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	defer serv.Close()

	// Wait forever
	for {
		time.Sleep(1 * time.Second)
	}
}
