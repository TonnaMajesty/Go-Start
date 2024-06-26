package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestTransmission(t *testing.T) {
	f, _ := os.ReadFile("/Users/tonnamajesty/Downloads/20230801-113220.jpeg")
	fs := base64.StdEncoding.EncodeToString(f)

	client := resty.New()
	resp, err := client.R().SetBody(map[string]interface{}{
		"fileData":    fs,
		"pointName":   "110kVall003_1",
		"msgReqId":    "20",
		"other":       "",
		"requestTime": 1689310544,
		"dayAndNight": "day",
	}).Post("http://172.20.40.21:9001/transmission-detection")
	fmt.Println(resp, err)
}
