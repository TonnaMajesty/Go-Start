package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.innoai.tech/component/svcutil/confhttp"
	"github.com/go-courier/envconf"
	"github.com/go-courier/httptransport/client"
	"github.com/go-courier/httptransport/transformers"

	"ai_center/client_ai_service"
)

var shmManger *ShmManager
var ctx context.Context
var aiClient *confhttp.ClientEndpoint

func init() {
	shmManger = NewShmManager("test-6", 100, 2*1024*1024)

	if err := shmManger.Open(); err != nil {
		fmt.Println(err)
		err := shmManger.Create()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	shmManger.Close()

	log.Fatal("")

	ctx = WithContext(context.Background())
	endpoint, _ := envconf.ParseEndpoint("http://localhost:4200")
	aiClient = &confhttp.ClientEndpoint{
		Endpoint: *endpoint,
		Timeout:  10 * time.Second,
	}
	aiClient.Init()
	aiClient.SetDefaults()
}

func main() {
	dir := "./test" // 指定文件夹路径

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".jpeg") {
				filePath := filepath.Join(dir, file.Name())

				// 打开文件
				f, err := os.Open(filePath)
				if err != nil {
					log.Printf("无法打开文件 %s: %s\n", filePath, err)
					continue
				}

				fb, err := io.ReadAll(f)
				if err != nil {
					log.Printf("读取文件错误 %s: %s\n", filePath, err)
					continue
				}

				fmt.Println(filePath)

				DoAnalyze(fb)
				time.Sleep(1 * time.Second)
				f.Close()
			}
		}
	}
}

var modelName = []string{
	"firesmog",
	"leave_station",
	"helmet",
	"gather",
	"fire_equipment",
}

var modelNameCount int

func getModelName() string {
	if modelNameCount+1 == len(modelName) {
		modelNameCount = 0
	} else {
		modelNameCount += 1
	}
	return modelName[modelNameCount]
}

func DoAnalyze(file []byte) {
	m := getModelName()
	m = "fire_equipment"
	req := client_ai_service.Analyze{
		ModelName: m,
		Product:   "charge",
	}

	binIndex, err := shmManger.WriteIntoShm(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileName := fmt.Sprintf("shm_%s_%d_%d_%d_%d", "test-6", 100, 2*1024*1024, binIndex, len(file))
	fmt.Println(fileName)
	fileHeader, err := BytesToFileHeaderWithName([]byte(fileName), fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Body.File = fileHeader
	req.Body.Params = `{"filter_params":{"MULTIPLEREPORT":[2],"REPORTONLY":[]},"smogfire_IOU_thres":0.05,"smogfire_det_thres":0.45}`
	req.Body.RequestID = fmt.Sprintf("firesmog-%d", time.Now().UnixMilli())
	resp, _, err := client_ai_service.NewClientAiService(aiClient, ctx).Analyze(&req)

	fmt.Println(m, resp, err)
}

func BytesToFileHeaderWithName(fileData []byte, filename string) (*multipart.FileHeader, error) {
	reader := bytes.NewReader(fileData)
	fileHeader, err := transformers.NewFileHeader(filename, filename, reader)
	if err != nil {
		return nil, err
	}

	return fileHeader, nil
}

var WithContext = confhttp.WithContextCompose(
	func(ctx context.Context) context.Context {
		return client.ContextWithDefaultHttpTransport(ctx, &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 0,
			}).DialContext,
			DisableKeepAlives:     true,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 30 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		})
	})
