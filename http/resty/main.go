package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	uuid "github.com/satori/go.uuid"
)

func main() {

	file, _ := os.ReadFile("/Users/tonnamajesty/Downloads/江苏点位.xlsx")

	resp, err := resty.New().R().
		SetHeaders(map[string]string{
			"Hello": "thank you",
		}).
		SetFileReader("excelFile", "test.xlsx", bytes.NewReader(file)).
		Post("http://localhost:4900/api/transmission-line/v0/import/jiangsu/point")

	fmt.Println(resp, err)

	//// 创建一个Resty客户端实例
	//client := resty.New()
	//
	//// 创建一个自定义的TLS配置
	//tlsConfig := &tls.Config{
	//	InsecureSkipVerify: true,
	//}
	//
	//// 将TLS配置设置到Resty客户端中
	//client.SetTLSClientConfig(tlsConfig)
	//
	//// 发起HTTPS请求
	//resp, err := client.R().Get("https://www.baidu.com/sugrec?&prod=pc_his&from=pc_web&json=1&sid=38516_36561_38687_38798_38767_38793_38844_38831_38582_38486_38801_38822_38839_38636_26350_22160&hisdata=%5B%7B%22time%22%3A1679482113%2C%22kw%22%3A%22golang%20%E5%AD%97%E7%AC%A6%E4%B8%B2%E8%BD%ACfloat64%22%7D%2C%7B%22time%22%3A1679552302%2C%22kw%22%3A%22waitgroup%22%7D%2C%7B%22time%22%3A1679640331%2C%22kw%22%3A%22golang%20case%20%E5%A4%9A%E4%B8%AA%22%7D%2C%7B%22time%22%3A1679654679%2C%22kw%22%3A%222006-01-02%2015%3A04%3A04%22%7D%2C%7B%22time%22%3A1679818583%2C%22kw%22%3A%22evicted%22%7D%2C%7B%22time%22%3A1679818662%2C%22kw%22%3A%22%E6%9F%A5%E7%9C%8B%E7%A3%81%E7%9B%98%E7%A9%BA%E9%97%B4%E5%A4%A7%E5%B0%8F%20linux%22%7D%2C%7B%22time%22%3A1679818801%2C%22kw%22%3A%22diskpressure%22%7D%2C%7B%22time%22%3A1679825310%2C%22kw%22%3A%22sda%20%E9%83%A8%E7%BD%B2%22%7D%2C%7B%22time%22%3A1679825315%2C%22kw%22%3A%22sda%20ssd%22%7D%2C%7B%22time%22%3A1679825373%2C%22kw%22%3A%22sda%20ssd%20%E4%BB%80%E4%B9%88%E5%8C%BA%E5%88%AB%22%7D%5D&_t=1686714504803&req=2&bs=%E7%99%BE%E5%BA%A6&csor=0")
	//
	//fmt.Println(resp, err)

	//client := resty.New()
	//client.SetHeaders(map[string]string{
	//	"Content-Type": "application/json",
	//})
	//client.R().SetBody(map[string]interface {
	//}{
	//	"msg_type": "text",
	//	"content": map[string]interface{}{
	//		"text": fmt.Sprintf("摄像头 %s 已离线, 杆塔号: %s, 电池电量: %d, 工作温度: %d, 请及时处理",
	//			"test", "test", 0, 0),
	//	},
	//}).Post("https://open.feishu.cn/open-apis/bot/v2/hook/6c1d1e08-ea28-43c7-a771-5a16de08363b")

	//client := resty.New()
	//client.SetBaseURL("https://srv-transmission-line-http---sentry.dev.innoai.tech")
	//client.SetAuthToken("eyJhbGciOiJSUzI1NiIsImtpZCI6IjgrbnRPT041U1VJIiwidHlwIjoiSldUIn0.eyJhdWQiOltdLCJleHAiOjE2ODkzMjM1NDEsImlhdCI6MTY4OTMxOTk0MSwiaXNzIjoiaHR0cDovL3Nzby5pbmR1c3RhaS5jb20iLCJqdGkiOiIxMTc1YzcwMS1iZGM3LTQ3ZjAtODFiNS0yOWUxMzg2OTUzMGMiLCJuYW1lIjoiYWRtaW44Iiwic2NvcGVzIjpbXSwic2VjdXJlQ29kZSI6IjE3Mi4yMC4zMC42OEA0d3lFbDQ5NzVsIiwic3ViIjoiMTQyNjE5NTc2NTMwMjM5NTI5OSIsInR5cGUiOiJhY2Nlc3NfdG9rZW4ifQ.I1b8Apud-RNw4MPw5X7041cHwdBhTBsXMqvOc0y7FigxaUrBeFg1o-Eoq-p98YuyZ41tjAaGV7GTctRT4iFhOO7Dq7h-iI30kAHc2Ca9vh-a07oYW8KrC8kiRroPAHni2FiiU5TC6pwUvAuHW7ithYbaCtAgLAzHA5GBblXC24n4aWt3OCyVhjOgtdcbNTOj8JOGP8xK1r08-8rqjaA25VEk4wwkwuKoDEBmlZR3uoboUXG5DSQPpC8DQbjUqem4jg_KYA-TuTCZKS2cb9eM_3CxByVMxtKDcV0oVyT6VXl60RXK9z_A5ciUWcu6eIGNIXgt7aRZp8Dq1z0LTGI-tg")
	//
	//res := Result{}
	//resp, err := client.R().SetQueryString("analyzeResult=ABNORMAL&analyzeResult=LOW_IMAGE_QUALITY&source=REAL_TIME_REPORT&size=4&offset=0").SetResult(&res).Get("/api/transmission-line/v0/inspection-image/list")
	//
	//if err != nil {
	//
	//}
	//client := resty.New()
	//resp, err := client.R().SetHeader("Content-Type", "application/json").SetBody(`name=1&hello=2`).Post("http://localhost:8080/")

	//fmt.Println(resp.String(), err)

	//client := resty.New()
	//tokenResp := TokenResp{}
	//resp, err := client.R().
	//	SetBody(
	//		fmt.Sprintf("grant_type=password&username=%s&password=%s&client_id=%s&client_secret=%s",
	//			"tzh", "123", "client", "secret")).
	//	SetResult(&tokenResp).
	//	Get("http://localhost:8080/auth/realms/product/protocol/openid-connect/token")
	//
	//fmt.Println(tokenResp, err)
	//fmt.Println(resp.String())

	requestID := uuid.NewV4()
	requestIDStr := requestID.String()
	fmt.Println(requestIDStr) // 3f675c44-cfdd-49ee-ae26-2bd283a0ba7e

	accessKey := "7b92579ef41449b18a06e33f9da56c44"
	//secretKey := "985e1b75c3f441afa96aebfb6a5d792"
	signTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(signTime)

	body := map[string]interface{}{
		"id": "1680761623004196864",
	}

	bodyStr, _ := json.Marshal(body)

	sign := string(bodyStr) + accessKey + requestIDStr + signTime
	content := sha256.Sum256([]byte(sign))
	signature := hex.EncodeToString(content[:])

	fmt.Println(signature)

}

type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type Result struct {
	Data  any `json:"data"`
	Total int `json:"total"`
}
