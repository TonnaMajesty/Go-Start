package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/xuri/excelize/v2"
)

func main() {
	// 定义命令行参数
	startTimePtr := flag.String("start", "2023-07-31 00:00:00", "Start time in the format '2006-01-02 15:04:05'")
	endTimePtr := flag.String("end", "2023-08-04 23:59:59", "End time in the format '2006-01-02 15:04:05'")
	flag.Parse()

	// 检查参数是否提供
	if *startTimePtr == "" || *endTimePtr == "" {
		fmt.Println("Usage: go run main.go -start <start_time> -end <end_time>")
		return
	}

	// 设置数据库连接信息
	db, err := sql.Open("postgres", "host=172.20.1.100 port=25434 user=root password=root@123456 dbname=ai_event_center sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 执行 SQL 查询
	rows, err := db.Query(`SELECT f_event_id 事件编号,CASE WHEN f_project_id=1731672592881811470 THEN '犀浦充电站' WHEN f_project_id=1736235039991029774 THEN '紫东正街充电站' WHEN f_project_id=1737403967865909262 THEN '柴家渡充电站' WHEN f_project_id=1737753042352095246 THEN '蓬溪客运充电站' WHEN f_project_id=1674357327764925056 THEN '麻石桥充电站' END 站点,f_device_name 车位,CASE WHEN f_scene_type='CHARGE_GUN_NOT_PLUGGED' THEN '充电枪未插回' WHEN f_scene_type='CHARGE_GUN_STATUS' THEN '充电枪状态（是否在桩上）' WHEN f_scene_type='VEHICLE_STOP' THEN '车辆停留' WHEN f_scene_type='PAT_EMERGENCY_BUTTON' THEN '乱拍紧急按钮' WHEN f_scene_type='FIRE_EQUIPMENT' THEN '消防器材' WHEN f_scene_type='POWER_DISTRIBUTION_INTRUSION' THEN '配电区单人进入' WHEN f_scene_type='INTRUSION' THEN '区域入侵' WHEN f_scene_type='ESCAPE_BILL' THEN '车辆逃单' WHEN f_scene_type='LICENSE_PLATE' THEN '车牌识别' WHEN f_scene_type='SMOKE_FIRE' THEN '烟火' END 模型,CASE WHEN f_type=1 THEN '异常' WHEN f_type=2 THEN '恢复' END 报警类型,to_char(TO_TIMESTAMP(f_receive_time) AT TIME ZONE 'asia/shanghai','YYYY-MM-dd hh24:MI:ss') 上报时间 FROM ai_event_center.t_event WHERE to_timestamp(f_receive_time) AT TIME ZONE 'asia/shanghai'>=to_timestamp($1,'yyyy-MM-dd HH24:MI:SS') AND to_timestamp(f_receive_time) AT TIME ZONE 'asia/shanghai'<=to_timestamp($2,'yyyy-MM-dd HH24:MI:SS') ORDER BY f_project_id,f_scene_type,f_receive_time`, *startTimePtr, *endTimePtr)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// 创建 Excel 文件
	xlsx := excelize.NewFile()
	sheetName := "Sheet1"
	xlsx.NewSheet(sheetName)

	// 设置表头
	xlsx.SetCellValue(sheetName, "A1", "事件编号")
	xlsx.SetCellValue(sheetName, "B1", "站点")
	xlsx.SetCellValue(sheetName, "C1", "车位")
	xlsx.SetCellValue(sheetName, "D1", "模型")
	xlsx.SetCellValue(sheetName, "E1", "报警类型")
	xlsx.SetCellValue(sheetName, "F1", "上报时间")

	// 将查询结果写入 Excel 文件
	row := 2
	for rows.Next() {
		var (
			eventID     string
			site        string
			parkingSpot string
			model       string
			alarmType   string
			reportTime  string
		)
		if err := rows.Scan(&eventID, &site, &parkingSpot, &model, &alarmType, &reportTime); err != nil {
			log.Fatal(err)
		}

		xlsx.SetCellValue(sheetName, fmt.Sprintf("A%d", row), eventID)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("B%d", row), site)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("C%d", row), parkingSpot)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("D%d", row), model)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("E%d", row), alarmType)
		xlsx.SetCellValue(sheetName, fmt.Sprintf("F%d", row), reportTime)

		row++
	}

	// 保存 Excel 文件
	err = xlsx.SaveAs("output.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Excel file successfully generated!")
}
