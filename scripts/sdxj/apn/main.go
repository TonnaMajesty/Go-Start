package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	client := NewSDXJClient("http://sdkshxj.scdl.cn:31200", "zhoup7477", "zhoup7477!", "sdxj", "fd013bca-5c26-432e-af90-c844bbaf6e47")

	codeDeviceMap := make(map[string]DeviceInfo)
	for _, code := range codelist {
		deviceList, err := client.GetDeviceInfo(code)
		if err != nil {
			fmt.Println("GetTerminalInfo Error", code, err)
			continue
		}

		for _, device := range deviceList {
			codeDeviceMap[device.DeviceID] = device
		}
	}

	codeIpMapStr, _ := json.Marshal(codeDeviceMap)
	os.WriteFile("./codeDeviceMap.json", codeIpMapStr, os.ModePerm)
}
