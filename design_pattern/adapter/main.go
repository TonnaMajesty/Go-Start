package main

import "fmt"

func main () {
	client := &client{}
	mac := &mac{}

	client.insertLightningConnectorIntoComputer(mac)

	windowsMachine := &windows{}
	windowsMachineAdapter := &windowsAdapter{
		windowMachine: windowsMachine,
	}

	client.insertLightningConnectorIntoComputer(windowsMachineAdapter)
}

// client 插入到computer上的接口
type client struct {

}

func (c *client) insertLightningConnectorIntoComputer (com computer) {
	fmt.Println("Client inserts Lightning connector into computer.")
	com.insertIntoLightningPort()
}

// computer interface
type computer interface {
	insertIntoLightningPort()
}

// mac 实现了 computer
type mac struct {}

func (m *mac) insertIntoLightningPort()  {
	fmt.Println("Lightning connector is plugged into mac machine.")
}

// win 没有实现computer
type windows struct {}

func (w *windows) insertIntoUSBPort() {
	fmt.Println("USB connector is plugged into windows machine.")
}

// windows 适配器
type windowsAdapter struct {
	windowMachine *windows
}

func (w *windowsAdapter) insertIntoLightningPort() {
	fmt.Println("Adapter converts Lightning signal to USB.")
	w.windowMachine.insertIntoUSBPort()
}

