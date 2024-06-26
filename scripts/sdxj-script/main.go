package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"git.innoai.tech/ai-apps/common/modelsx"
	"github.com/go-courier/sqlx/v2/datatypes"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

var (
	// 1秒钟3个
	secondLimit = rate.NewLimiter(Per(3, 1*time.Second), 1)
	//60秒180个
	minuteLimit = rate.NewLimiter(Per(180, 1*time.Minute), 20)
	limiter     = MultiLimiter(secondLimit, minuteLimit)
	manager     = NewManger(10)
)

func main() {
	var start string
	var end string
	var codelist string
	var district string

	var rootCmd = &cobra.Command{
		Use:   "transmission",
		Short: "CLI for download transmission picture",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Start: %s\n", start)
			fmt.Printf("End: %s\n", end)
			fmt.Printf("Code: %s\n", codelist)
			codel := []string{}
			if codelist != "" {
				codel = strings.Split(codelist, ",")
			}
			ctx := context.Background()
			//manager.Run(ctx)
			Process(ctx, start, end, codel, district)

			// 定时监测是否完成
			//ticker := time.NewTicker(10 * time.Second)
			//for {
			//	select {
			//	case <-ticker.C:
			//		if len(manager.MImageChan) == 0 {
			//			return
			//		}
			//	}
			//}
		},
	}

	rootCmd.Flags().StringVarP(&start, "start", "s", "", "Start parameter")
	rootCmd.Flags().StringVarP(&end, "end", "e", "", "End parameter")
	rootCmd.Flags().StringVarP(&codelist, "code", "c", "", "code parameter, split by ,")
	rootCmd.Flags().StringVarP(&district, "district", "d", "", "district short, eg: cd,xc,gwcgy,gwcd")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
}

func Process(ctx context.Context, start, end string, codelist []string, district string) {
	startTimestamp, err := datatypes.ParseTimestampFromStringWithLayout(start, "2006-01-02 15:04:05")
	if err != nil {
		logrus.Errorf("ParseTimestampFromStringWithLayout startTimestamp %s failed[err:%s]", start, err.Error())
		return
	}

	endTimestamp, err := datatypes.ParseTimestampFromStringWithLayout(end, "2006-01-02 15:04:05")
	if err != nil {
		logrus.Errorf("ParseTimestampFromStringWithLayout endTimestamp %s failed[err:%s]", end, err.Error())
		return
	}

	if district != "" {
		if l, ok := districtMap[district]; ok {
			codelist = l
		} else {
			logrus.Errorf("district not exists")
			return
		}
	}

	//client := NewSDXJClient("http://sdkshxj.scdl.cn:31200", "sdxj-siji", "SdXj-siji", "sdxj", "fd013bca-5c26-432e-af90-c844bbaf6e47")
	client := NewSDXJClient("http://sdkshxj.scdl.cn:31200", "zhoup7477", "zhoup7477!", "sdxj", "fd013bca-5c26-432e-af90-c844bbaf6e47")
	codeDeviceInfoMap := make(map[string][]DeviceInfo)

	if _, err = os.Stat("./codeDeviceInfoMap.json"); err == nil {
		codeDeviceInfoMapByte, err := os.ReadFile("./codeDeviceInfoMap.json")
		if err == nil {
			err = json.Unmarshal(codeDeviceInfoMapByte, &codeDeviceInfoMap)
			if err != nil {
				fmt.Println("Unmarshal codeDeviceInfoMap error", err)
			}
		}
	}

	if len(codeDeviceInfoMap) == 0 || len(codelist) > 0 {
		if len(codelist) == 0 {
			codelist = append(CodeListCD, CodeListXC...)
		}
		newCodeDeviceInfoMap := make(map[string][]DeviceInfo)

		for _, code := range codelist {
			if i, ok := codeDeviceInfoMap[code]; ok {
				newCodeDeviceInfoMap[code] = i
			} else {
				deviceInfoList, err := client.GetDeviceInfo(code)
				if err != nil {
					logrus.Errorf("GetDeviceInfo %s failed[err:%s]", code, err.Error())
					continue
				}
				newCodeDeviceInfoMap[code] = deviceInfoList
			}
		}

		codeDeviceInfoMap = newCodeDeviceInfoMap
		// 保存到本地
		codeDeviceInfoMapStr, _ := json.Marshal(codeDeviceInfoMap)
		os.WriteFile("./codeDeviceInfoMap.json", codeDeviceInfoMapStr, os.ModePerm)
	}

	var wg sync.WaitGroup

	p, _ := ants.NewPool(1)
	defer ants.Release()

	for code, deviceInfoList := range codeDeviceInfoMap {
		//deviceStateInfo, err := client.GetDeviceStateInfo(deviceMN)
		//if err != nil {
		//	logrus.Errorf("SDXJCron GetDeviceStateInfo code %s devicemn %s failed[err:%s]", code, deviceMN, err.Error())
		//	continue
		//}
		//logrus.Infof("SDXJCron GetDeviceStateInfo code %s devicemn %s success %+v", code, deviceMN, deviceStateInfo)

		deviceInfoList := deviceInfoList
		code := code
		wg.Add(1)
		_ = p.Submit(func() {
			GetAndDownloadPicture(ctx, client, deviceInfoList, startTimestamp, endTimestamp, code)
			wg.Done()
		})
	}
	wg.Wait()
}

func GetAndDownloadPicture(ctx context.Context, client *SDXJClient, deviceInfoList []DeviceInfo, startTimestamp, endTimestamp modelsx.Timestamp, code string) {
	for _, deviceInfo := range deviceInfoList {
		if deviceInfo.ChannelID == "-1" {
			continue
		}
		pictureListResult, err := client.GetHistoryPicture(deviceInfo.DeviceID, startTimestamp, endTimestamp, 1, 50)
		if err != nil {
			logrus.Errorf("GetHistoryPicture %s_%s failed[err:%s]", code, deviceInfo.ChannelID, err.Error())
			return
		}

		fmt.Printf("pictureListResult %s_%s total %d pages %d pageNum %d pageSize %d\n", code, deviceInfo.ChannelID, pictureListResult.Total, pictureListResult.Pages, pictureListResult.PageNum, pictureListResult.PageSize)

		if pictureListResult.Total == 0 {
			logrus.Infof("GetHistoryPicture %s_%s total 0", code, deviceInfo.ChannelID)
			return
		}

		for i := 1; i <= (pictureListResult.Total/50)+1; i++ {
			pictureListResult, err := client.GetHistoryPicture(deviceInfo.DeviceID, startTimestamp, endTimestamp, i, 50)
			if err != nil {
				logrus.Errorf("GetHistoryPicture %s_%s failed[err:%s]", code, deviceInfo.ChannelID, err.Error())
				continue
			}

			pictureList := pictureListResult.PictureList

			for _, picture := range pictureList {

				//manager.AddJob(&Pic{
				//	url:        picture.Url,
				//	uploadTime: picture.UploadTime,
				//	code:       code,
				//	channelID:  deviceInfo.ChannelID,
				//	towerName:  picture.TowerName,
				//})

				// 下载图片
				logrus.Infof("start downloadPicture %s", picture.Url)
				content, err := downloadPicture(ctx, fmt.Sprintf("%s%s", "http://sdkshxj.scdl.cn:31200", picture.Url))
				if err != nil {
					logrus.Errorf("downloadPicture %s failed[err:%s]", picture.Url, err.Error())
					continue
				}

				timestamp, err := datatypes.ParseTimestampFromStringWithLayout(picture.UploadTime, "2006-01-02 15:04:05")
				if err != nil {
					logrus.Errorf("SDXJCron ParseTimestampFromStringWithLayout %s failed[err:%s]", picture.UploadTime, err.Error())
					continue
				}

				subDir := fmt.Sprintf("./%s_%s/image", code, deviceInfo.ChannelID)
				_, err = os.Stat(subDir)
				if err != nil {
					if os.IsNotExist(err) {
						err := os.MkdirAll(subDir, fs.ModePerm)
						if err != nil {
							logrus.Errorf("MkdirAll %s failed[err:%s]", subDir, err.Error())
							continue
						}
					}
				}

				// 获取图片后缀
				suffix := strings.Split(picture.Url, ".")[len(strings.Split(picture.Url, "."))-1]
				if len(strings.Split(picture.Url, ".")) == 1 {
					suffix = "jpg"
				}

				// 保存图片到本地
				poutfile := fmt.Sprintf("%s/%s_%s.%s", subDir, picture.TowerName, timestamp.Format("20060102150405"), suffix)
				os.WriteFile(poutfile, content, os.ModePerm)
			}
		}
	}
}

type Pic struct {
	url        string
	uploadTime string
	code       string
	channelID  string
	towerName  string
}

func DealPicture(ctx context.Context, p *Pic) {
	logrus.Infof("start downloadPicture %s", p.url)
	content, err := downloadPicture(ctx, fmt.Sprintf("%s%s", "http://sdkshxj.scdl.cn:31200", p.url))
	if err != nil {
		logrus.Errorf("downloadPicture %s failed[err:%s]", p.url, err.Error())
		return
	}

	timestamp, err := datatypes.ParseTimestampFromStringWithLayout(p.uploadTime, "2006-01-02 15:04:05")
	if err != nil {
		logrus.Errorf("SDXJCron ParseTimestampFromStringWithLayout %s failed[err:%s]", p.uploadTime, err.Error())
		return
	}

	subDir := fmt.Sprintf("./%s_%s/image", p.code, p.channelID)
	_, err = os.Stat(subDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(subDir, fs.ModePerm)
			if err != nil {
				logrus.Errorf("MkdirAll %s failed[err:%s]", subDir, err.Error())
				return
			}
		}
	}

	// 获取图片后缀
	suffix := strings.Split(p.url, ".")[len(strings.Split(p.url, "."))-1]
	if len(strings.Split(p.url, ".")) == 1 {
		suffix = "jpg"
	}

	// 保存图片到本地
	poutfile := fmt.Sprintf("%s/%s_%s.%s", subDir, p.towerName, timestamp.Format("20060102150405"), suffix)
	os.WriteFile(poutfile, content, os.ModePerm)
}

func downloadPicture(ctx context.Context, url string) ([]byte, error) {
	if err := limiter.Wait(ctx); err == nil {
		resp, err := (&http.Client{
			Timeout: 10 * time.Second,
		}).Get(url)
		if err != nil {
			return nil, err
		}
		defer func() {
			resp.Body.Close()
		}()

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusTooManyRequests {
				fmt.Println("downloadPicture failed with StatusTooManyRequests")
				// todo 限流重试
			}
			return nil, fmt.Errorf("downloadPicture unexpected status code: %d", resp.StatusCode)
		}

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		return content, nil
	}

	return nil, nil
}
