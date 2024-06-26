package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	// 读取 Excel 文件
	xlsx, err := excelize.OpenFile("/Users/tonnamajesty/Downloads/点位清单1(1).xlsx")
	if err != nil {
		log.Fatal(err)
	}

	// 定义列表用于保存数据
	var dataList []Data
	// 定义 map 用于保存数据
	dataMap := make(map[string]Data)

	// 解析 Excel 文件，将数据保存到列表中
	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	for i, row := range rows {
		if i == 0 {
			continue // 跳过标题行
		}
		if len(row) < 2 {
			continue // 跳过不完整的行
		}

		data := parsePointName(row[0])
		data.Code = row[1]

		// 将电压等级、线路名称、杆塔号和安装位置作为 key，将记录存储在 map 中
		key := fmt.Sprintf("%s%s%s_%d", data.Voltage, data.Line, data.Tower, data.Position)
		if _, ok := dataMap[key]; !ok {
			dataMap[key] = data
		} else {
			if dataMap[key].Code != data.Code {
				// 如果 key 已经存在，但是点位编码不同，打印不同的点位编码
				fmt.Printf("Different code: %s, %s\n", dataMap[key].Code, data.Code)
			}
			// 如果 key 已经存在，打印重复的key
			fmt.Printf("Duplicate key: %s\n", key)
		}

		dataList = append(dataList, data)
	}

	// 打印解析后的数据
	for _, data := range dataMap {
		if data.Code == "" {
			continue
		}
		// 判断是否以INDUSTAI开头
		var cameraType string
		if strings.HasPrefix(data.Code, "INDUSTAI") {
			cameraType = "enums.CAMERA_TYPE__REAL_TIME"
		} else {
			cameraType = "enums.CAMERA_TYPE__CYCLE"
		}
		//fmt.Printf("Voltage: %s, Line: %s, Tower: %s, Position: %d, Code: %s\n", data.Voltage, data.Line, data.Tower, data.Position, data.Code)

		fmt.Printf(`	{
		CameraType:    %s,
		Platform:      enums.CAMERA_PLATFORM__SDXJ,
		MachineNo:     "%s",
		DeviceMN:      "",
		OperationTeam: "",
		VoltageLevel:  "%s",
		Line:          "%s",
		Tower:         "%s",
		CustomID:      %d,
		Location:      custom.INSTALL_LOCATION_UNKNOWN,
		LocationInfo:  "",
		ExtraInfo: models.CameraExtraInfo{
			IsSupportControl: false,
			TerminalCode:     "%s",
		},
	},`, cameraType, data.Code, data.Voltage, data.Line, data.Tower, data.Position, data.Code)
	}
}

// 定义数据结构
type Data struct {
	Voltage  string
	Line     string
	Tower    string
	Position int
	Code     string
}

// 解析点位名称
func parsePointName(pointName string) Data {
	data := Data{}

	// 提取电压等级
	voltageRegex := regexp.MustCompile(`^(110|220)kV`)
	match := voltageRegex.FindString(pointName)
	if match != "" {
		data.Voltage = match
	}

	// 提取线路名称
	lineRegex := regexp.MustCompile(`kV(.+?线)\d`)
	matchs := lineRegex.FindStringSubmatch(pointName)
	if len(matchs) >= 2 {
		data.Line = strings.TrimSpace(matchs[1])
	}
	// 提取安装位置
	// 提取杆塔号和安装位置
	towerPositionRegex := regexp.MustCompile(`线(\d+)号?(?:_(\d))?(?:_(\d))?$`)
	matchs = towerPositionRegex.FindStringSubmatch(pointName)
	if len(matchs) >= 2 {
		data.Tower = matchs[1]
		if len(match) >= 3 && matchs[2] != "" {
			position, _ := strconv.Atoi(matchs[2])
			data.Position = position
		}
	}

	return data
}
