package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/go-courier/sqlx/v2/datatypes"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	var start string
	var end string
	var codelist string
	var district string

	var rootCmd = &cobra.Command{
		Use:   "transmission-stat",
		Short: "CLI for  transmission stat",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Start: %s\n", start)
			fmt.Printf("End: %s\n", end)
			fmt.Printf("Code: %s\n", codelist)
			codel := []string{}
			if codelist != "" {
				codel = strings.Split(codelist, ",")
			}
			ctx := context.Background()
			//start = "2023-10-26 22:50:07"
			//end = "2023-11-30 16:20:22"
			Process(ctx, start, end, codel, district)
		},
	}

	rootCmd.Flags().StringVarP(&start, "start", "s", "", "Start parameter")
	rootCmd.Flags().StringVarP(&end, "end", "e", "", "End parameter")
	rootCmd.Flags().StringVarP(&codelist, "code", "c", "", "code parameter, split by ,")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
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

	if endTimestamp.Unix()-startTimestamp.Unix() > 24*60*60*31 {
		logrus.Errorf("时间段不能超过一个月")
		return
	}

	// 设置数据库连接信息
	db, err := sql.Open("postgres", "host=10.176.200.105 port=25434 user=root password=Ai@indust.PG!231206 dbname=transmission_line sslmode=disable")
	//db, err := sql.Open("postgres", "host=172.25.1.147 port=5432 user=user_rw password=Abc123456 dbname=transmission_line sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ariealPass, err := getCount(db, 3, 2, startTimestamp.Unix(), endTimestamp.Unix())
	if err != nil {
		return
	}
	ariealReject, err := getCount(db, 3, 3, startTimestamp.Unix(), endTimestamp.Unix())
	if err != nil {
		return
	}
	smokePass, err := getCount(db, 7, 2, startTimestamp.Unix(), endTimestamp.Unix())
	if err != nil {
		return
	}
	smokeReject, err := getCount(db, 7, 3, startTimestamp.Unix(), endTimestamp.Unix())
	if err != nil {
		return
	}

	ariealPassMap := make(map[string]int)
	ariealRejectMap := make(map[string]int)
	smokePassMap := make(map[string]int)
	smokeRejectMap := make(map[string]int)

	for _, s := range ariealPass {
		ariealPassMap[fmt.Sprintf("%s_%s", s.DeviceCode, s.ChannelID)] = s.Count
	}

	for _, s := range ariealReject {
		ariealRejectMap[fmt.Sprintf("%s_%s", s.DeviceCode, s.ChannelID)] = s.Count
	}

	for _, s := range smokePass {
		smokePassMap[fmt.Sprintf("%s_%s", s.DeviceCode, s.ChannelID)] = s.Count
	}

	for _, s := range smokeReject {
		smokeRejectMap[fmt.Sprintf("%s_%s", s.DeviceCode, s.ChannelID)] = s.Count
	}

	var ariealMiss []miss
	for s, i := range ariealRejectMap {
		mis := float64(i) / float64(ariealPassMap[s]+i)
		ariealMiss = append(ariealMiss, miss{
			Dev:     s,
			MisRate: mis,
		})
	}

	var smokeMiss []miss
	for s, i := range smokeRejectMap {
		mis := float64(i) / float64(smokePassMap[s]+i)
		smokeMiss = append(smokeMiss, miss{
			Dev:     s,
			MisRate: mis,
		})
	}

	sort.Slice(ariealMiss, func(i, j int) bool {
		return ariealMiss[i].MisRate > ariealMiss[j].MisRate
	})

	sort.Slice(smokeMiss, func(i, j int) bool {
		return smokeMiss[i].MisRate > smokeMiss[j].MisRate
	})

	resultStr := "【高空异物】\n"

	for _, m := range ariealMiss {
		resultStr += fmt.Sprintf("%s %.2f\n", m.Dev, m.MisRate)
	}

	resultStr += "【烟火】\n"

	for _, m := range smokeMiss {
		resultStr += fmt.Sprintf("%s %.2f\n", m.Dev, m.MisRate)
	}
	filePath := fmt.Sprintf("./transmission-stat-%s.txt", startTimestamp.Format("20060102150405"))

	err = os.WriteFile(filePath, []byte(resultStr), os.ModePerm)
	if err != nil {
		logrus.Errorf("Write Result err %s", err)
		return
	}

	// 获取每个点位每天的通过的数量

	ariealPassCount, err := getPassCountByPoint(db, 3, startTimestamp.Unix(), endTimestamp.Unix())
	if err != nil {
		return
	}

	passResultStr := "【高空异物】\n"
	for _, m := range ariealPassCount {
		passResultStr += fmt.Sprintf("%s    %s    %d\n", m.LineName, m.TowerName, m.Count)
	}

	passResultStr += "【烟火】\n"

	smokePassCount, err := getPassCountByPoint(db, 7, startTimestamp.Unix(), endTimestamp.Unix())
	for _, m := range smokePassCount {
		passResultStr += fmt.Sprintf("%s    %s    %d\n", m.LineName, m.TowerName, m.Count)
	}

	passFilePath := fmt.Sprintf("./transmission-pass-stat-%s.txt", startTimestamp.Format("20060102150405"))

	err = os.WriteFile(passFilePath, []byte(passResultStr), os.ModePerm)
	if err != nil {
		logrus.Errorf("Write Result err %s", err)
		return
	}
}

type stat struct {
	DeviceCode string
	ChannelID  string
	Count      int
}

type miss struct {
	Dev     string
	MisRate float64
}

type passStat struct {
	LineName  string
	TowerName string
	Count     int
}

func getPassCountByPoint(db *sql.DB, analyzeResult int, start, end int64) ([]passStat, error) {
	sqlStr := `select f_line, f_tower, count(*) from t_inspection_image_analyze where f_analyze_result = $1 and f_audit_result = 2 and f_inspected_at >= $2 and f_inspected_at <= $3 group by f_line, f_tower`
	rows, err := db.Query(sqlStr, analyzeResult, start, end)
	if err != nil {
		logrus.Errorf("query failed %s", err)
		return nil, err
	}

	var result []passStat
	for rows.Next() {
		// 读取每一行数据
		var (
			line  string
			tower string
			count int
		)

		err := rows.Scan(&line, &tower, &count)
		if err != nil {
			logrus.Errorf("scan failed %s", err)
			return nil, err
		}

		result = append(result, passStat{
			LineName:  line,
			TowerName: tower,
			Count:     count,
		})
	}

	rows.Close()

	return result, nil
}

func getCount(db *sql.DB, analyzeResult, auditResult int, start, end int64) ([]stat, error) {
	sqlStr := `select f_device_code, f_channel_id, count(*) from t_inspection_image_analyze where f_analyze_result = $1 and f_audit_result = $2 and f_inspected_at >= $3 and f_inspected_at <= $4 group by f_device_code, f_channel_id`
	rows, err := db.Query(sqlStr, analyzeResult, auditResult, start, end)
	if err != nil {
		logrus.Errorf("query failed %s", err)
		return nil, err
	}

	var result []stat
	for rows.Next() {
		// 读取每一行数据
		var (
			deviceCode string
			channelID  string
			count      int
		)

		err := rows.Scan(&deviceCode, &channelID, &count)
		if err != nil {
			logrus.Errorf("scan failed %s", err)
			return nil, err
		}

		result = append(result, stat{
			DeviceCode: deviceCode,
			ChannelID:  channelID,
			Count:      count,
		})
	}

	rows.Close()

	return result, nil
}
