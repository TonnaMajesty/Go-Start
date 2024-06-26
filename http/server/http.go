package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

func main() {
	var count int64
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&count, 1)
		fmt.Printf("==============================, %d\n", count)
		// 打印请求URL
		fmt.Printf("Request URL: %s\n", r.URL.Path)

		// 打印请求Header
		fmt.Println("Request Headers:")
		for key, value := range r.Header {
			fmt.Printf("%s: %s\n", key, value)
		}

		// 读取请求Body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		// 打印请求Body
		fmt.Printf("Request Body: %s\n", string(body))

		// 返回响应

		if r.URL.Path == "/auth/realms/product/protocol/openid-connect/token" {
			res := map[string]interface{}{
				"access_token": "kvakshfjsdafh3edfjkasvjkashgjsdafhasf",
				"expires_in":   3600,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)

		} else if r.URL.Path == "/sdxj/device/imagedev_signal_battery_info" {
			res := map[string]interface{}{
				"successful": true,
				"resultValue": map[string]interface{}{
					"onlineState":    1,
					"remainElectric": 0.6,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)
		} else if r.URL.Path == "/tDeviceInfo" {
			res := map[string]interface{}{
				"successful": true,
				"resultValue": map[string]interface{}{
					"itemCount": 1,
					"items": []DeviceInfo{
						{
							DeviceID:    strconv.FormatInt(time.Now().UnixMilli(), 10),
							LineName:    "线路1",
							TowerName:   "杆塔1",
							DeviceState: 1,
							UpdateTime:  time.Now().UnixMilli(),
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)
		} else if r.URL.Path == "/race" {
			req := map[string]interface{}{}
			json.Unmarshal(body, &req)
			fmt.Println(req)
			res := map[string]interface{}{
				"resId":         "响应ID，与请求ID相同",
				"resCode":       "0", // 成功
				"resDesc":       "分析结果描述",
				"resultJsonStr": `{"result":[{"tag":"fire","score":0.87,"bndbox":[52,648,253,952]}]}`,
				"resTime":       1693300824603,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)
		} else if r.URL.Path == "/callback" {
			req := map[string]interface{}{}
			json.Unmarshal(body, &req)
			fmt.Println(req)

			res := map[string]interface{}{
				"code": rand.Uint64(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)
		} else {
			//if rand.Uint64()%20 > 9 {
			//	time.Sleep(2 * time.Second)
			//}

			//time.Sleep(1100 * time.Millisecond)
			res := map[string]interface{}{
				"code": 200000,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)
		}
	})

	// 启动HTTP Server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type DeviceInfo struct {
	DeviceID    string `json:"id"`           // 设备ID
	LineName    string `json:"line_Name"`    // 线路名称
	TowerName   string `json:"tower_Name"`   // 杆塔名称
	DeviceState int    `json:"device_State"` // 设备状态 0：离线 1：在线
	UpdateTime  int64  `json:"update_Time"`  // 最新同步时间
}
