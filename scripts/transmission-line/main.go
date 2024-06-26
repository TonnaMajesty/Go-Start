package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"git.innoai.tech/component/confs3/v2"
	datatypes2 "git.innoai.tech/component/datatypes"
	"github.com/go-courier/sqlx/v2/datatypes"
	"github.com/panjf2000/ants/v2"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	//s3 = &confs3.ObjectDB{
	//	Endpoint:        "obs.cn-north-4.myhuaweicloud.com",
	//	AccessKeyID:     "7SJ437WMKSAME1SQERCS",
	//	SecretAccessKey: "bK3nHhrwdiL9oCoTxbR6nxvJMZMKdmsbYZvekEUv",
	//	BucketName:      "dev-industai-app",
	//	Secure:          false,
	//}

	s3 = &confs3.ObjectDB{
		Endpoint:        "10.176.200.105:39000",
		AccessKeyID:     "ACh5Jnh4CTxeg9DK3BAN33Pt",
		SecretAccessKey: "AKJKSARUuMarPPcVE5KZBX54TB79bMBJdR2E",
		BucketName:      "dev-industai-app",
		Secure:          false,
	}

	// 1秒钟3个
	secondLimit = rate.NewLimiter(Per(3, 1*time.Second), 1)
	//60秒180个
	minuteLimit = rate.NewLimiter(Per(180, 1*time.Minute), 20)
	limiter     = MultiLimiter(secondLimit, minuteLimit)
)

func main() {
	var start string
	var end string
	var codelist string
	var district string
	var source string

	var rootCmd = &cobra.Command{
		Use:   "transmission",
		Short: "CLI for download transmission picture",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Start: %s\n", start)
			fmt.Printf("End: %s\n", end)
			fmt.Printf("Code: %s\n", codelist)
			fmt.Printf("District: %s\n", district)
			fmt.Printf("Source: %s\n", source)
			codel := []string{}
			if codelist != "" {
				codel = strings.Split(codelist, ",")
			}
			var isRace bool
			if source == "race" {
				isRace = true
			}
			ctx := context.Background()
			//start = "2023-10-26 22:50:07"
			//end = "2023-10-30 16:20:22"

			Process(ctx, start, end, codel, district, isRace)
		},
	}

	rootCmd.Flags().StringVarP(&start, "start", "s", "", "Start parameter")
	rootCmd.Flags().StringVarP(&end, "end", "e", "", "End parameter")
	rootCmd.Flags().StringVarP(&codelist, "code", "c", "", "code parameter, split by ,")
	rootCmd.Flags().StringVarP(&codelist, "source", "r", "", "data source")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Process(ctx context.Context, start, end string, codelist []string, district string, isRace bool) {
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

	if endTimestamp.Unix()-startTimestamp.Unix() > 24*60*60 {
		logrus.Errorf("时间段不能超过一天")
		return
	}

	// 设置数据库连接信息
	db, err := sql.Open("postgres", "host=10.176.200.105 port=25434 user=root password=Ai@indust.PG!231206 dbname=transmission_line sslmode=disable")
	//db, err := sql.Open("postgres", "host=172.25.1.147 port=5432 user=user_rw password=Abc123456 dbname=transmission_line sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var (
		limit  = 5000
		offset = 0
	)

	var wg sync.WaitGroup
	p, _ := ants.NewPool(10)
	defer ants.Release()

	for {
		var sqlStr string
		if isRace {
			sqlStr = `select f_image_analyze_id,f_alg_type,f_object_url,f_analyze_result,f_extra,f_device_code,f_channel_id,f_audit_result,f_inspected_at 
from t_inspection_image_analyze where f_inspected_at > $1 and f_inspected_at < $2 and f_analyze_state = 3 and f_analyze_result in (2,3,7) and f_callback_extra like '%isRace":true%' limit $3 offset $4`
		} else {
			sqlStr = `select f_image_analyze_id,f_alg_type,f_object_url,f_analyze_result,f_extra,f_device_code,f_channel_id,f_audit_result,f_inspected_at 
from t_inspection_image_analyze where f_inspected_at > $1 and f_inspected_at < $2 and f_analyze_state = 3 and f_analyze_result in (2,3,7) and f_callback_extra not like '%isRace":true%' limit $3 offset $4`
		}

		fmt.Println(sqlStr)
		rows, err := db.Query(sqlStr, startTimestamp.Unix(), endTimestamp.Unix(), limit, offset)
		if err != nil {
			logrus.Errorf("query failed %s", err)
			return
		}

		// 遍历查询结果
		hasData := false
		for rows.Next() {
			hasData = true
			// 读取每一行数据
			var (
				imageAnalyzeId int
				algType        int
				objectUrl      string
				analyzeResult  int
				extra          string
				deviceCode     string
				channelID      string
				auditResult    int
				inspectedAt    int
			)

			err := rows.Scan(&imageAnalyzeId, &algType, &objectUrl, &analyzeResult, &extra, &deviceCode, &channelID, &auditResult, &inspectedAt)
			if err != nil {
				log.Fatal(err)
			}

			// 处理数据
			fmt.Println(imageAnalyzeId, algType, objectUrl, analyzeResult, extra, deviceCode, channelID, auditResult, inspectedAt)

			wg.Add(1)
			_ = p.Submit(func() {
				DealImage(ctx, imageAnalyzeId, algType, objectUrl, analyzeResult, extra, deviceCode, channelID, auditResult, inspectedAt)
				wg.Done()
			})
		}

		rows.Close()

		// 更新偏移量
		offset += limit
		if !hasData {
			break
		}
	}

	wg.Wait()
}

func DealImage(ctx context.Context, imageAnalyzeId int, algType int, objectUrl string, analyzeResult int, extra string, deviceCode string, channelID string, auditResult int, inspectedAt int) {
	var subDir string
	if algType == 1 {
		subDir = fmt.Sprintf("./%s_%s/高空", deviceCode, channelID)
	} else {
		subDir = fmt.Sprintf("./%s_%s/烟火", deviceCode, channelID)
	}

	if analyzeResult == 2 { // 正常
		subDir = subDir + "-正常"

	} else if analyzeResult == 3 || analyzeResult == 7 { // 高空异物\烟火
		subDir = subDir + "异常"
		if auditResult == 2 { // 通过
			subDir = subDir + "-审核通过"
		} else if auditResult == 3 { // 未通过
			subDir = subDir + "-审核不通过"
		} else {
			return
		}
	} else {
		return
	}

	_, err := os.Stat(subDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(subDir, fs.ModePerm)
			if err != nil {
				logrus.Errorf("MkdirAll %s failed[err:%s]", subDir, err.Error())
				return
			}
		}
	}

	dUrl, err := GetDownloadURL(ctx, objectUrl)
	if err != nil {
		logrus.Errorf("GetDownloadURL %s err %s", dUrl, err)
		return
	}

	if err := limiter.Wait(ctx); err == nil {
		imgResp, err := http.Get(dUrl)
		if err != nil {
			logrus.Errorf("http.Get %s err %s", dUrl, err)
			return
		}

		if imgResp.StatusCode != http.StatusOK {
			if imgResp.StatusCode == http.StatusTooManyRequests {
				logrus.Errorf("downloadPicture failed with StatusTooManyRequests")
			}
			return
		}

		defer imgResp.Body.Close()

		imgByte, err := io.ReadAll(imgResp.Body)
		if err != nil {
			logrus.Errorf("ReadAll %s err %s", dUrl, err)
			return
		}

		// 获取图片后缀
		fmt.Println(strings.Split(objectUrl, "."))
		suffix := strings.Split(objectUrl, ".")[len(strings.Split(objectUrl, "."))-1]
		if len(strings.Split(objectUrl, ".")) == 1 {
			suffix = "jpg"
		}

		// 保存图片到本地
		poutfile := fmt.Sprintf("%s/%d-%s.%s", subDir, imageAnalyzeId, time.Unix(int64(inspectedAt), 0).Format("20060102150405"), suffix)
		err = os.WriteFile(poutfile, imgByte, os.ModePerm)
		if err != nil {
			logrus.Errorf("Write Image err %s", err)
			return
		}

		// 保存分析结果
		resultFile := fmt.Sprintf("%s/%d.json", subDir, imageAnalyzeId)
		err = os.WriteFile(resultFile, []byte(extra), os.ModePerm)
		if err != nil {
			logrus.Errorf("Write Result err %s", err)
			return
		}
	}
}

func GetDownloadURL(ctx context.Context, fileIdentification string) (string, error) {
	//prefix := "https://obs.sc-region-1.sgic.sgcc.com.cn/sdxlkshznxsxt-obs001/cloud/i1/sdxj-store"
	prefix := "http://sdkshxj.scdl.cn:31200"
	// 华雁系统图片路径
	if strings.HasPrefix(strings.TrimLeft(fileIdentification, "/"), "sdxj-store") {
		//return prefix + fileIdentification[strings.LastIndex(fileIdentification, "sdxj-store")+10:], nil
		return prefix + fileIdentification, nil
	}

	// minio地址 asset://xxx
	if strings.HasPrefix(fileIdentification, "asset") {
		address, err := datatypes2.ParseAddress(fileIdentification)
		if err != nil {
			return "", errors.Wrap(err, "GetDownloadURL ParseAddress Error")
		}

		meta, err := s3.StatsObject(ctx, *address)
		if err != nil {
			return "", errors.Wrap(err, "GetDownloadURL StatsObject Error")
		}

		url, err := s3.ProtectURL(ctx, meta, time.Duration(500)*time.Second)
		if err != nil {
			return "", errors.Wrap(err, "GetDownloadURL ProtectURL Error")
		}

		return url.String(), nil
	}

	return "", errors.New("fileIdentification illegal")
}
