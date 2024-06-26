package main

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/xuri/excelize/v2"

	_ "time/tzdata"

	"github.com/go-courier/sqlx/v2/builder"
	sqlxDatatypes "github.com/go-courier/sqlx/v2/datatypes"
)

const (
	apiURL     = "http://10.168.4.218/api/gas-station-gateway/v0/statistic/0/vehicle/model/export?authorization=Bearer%20eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SUQiOiIxNTg3OTA1Nzc4MDQ2NDc2NjIyIiwiZXhwIjoxNjkxMDU1MTIxLCJpbnRlcmZhY2VzIjpbIkhhbmRsZUV2ZW50IiwiQ291bnREZXZpY2UiLCJGaWx0ZXJBbmFseXNpc1R5cGUiLCJGaWx0ZXJEZXZpY2VMaXN0IiwiR2V0T2JqZWN0VXJpIiwiTGlzdEV2ZW50IiwiQ291bnRFdmVudCIsIlBlcmlvZGljU3RhdGlzdGljIiwiVmVoaWNsZURpc3RyaWJ1dGlvblN0YXRpc3RpYyIsIldlYnNvY2tldENvbm5lY3Rpb24iLCJHZXRPYmplY3RDb250ZW50IiwiVmVoaWNsZVBvc2l0aW9uU3RhdGlzdGljIiwiVmVoaWNsZU1vZGVsRGlzdHJpYnV0aW9uU3RhdGlzdGljIiwiVmVoaWNsZU1vZGVsU3RhdGlzdGljIiwiVmVoaWNsZU1vZGVsU3RhdGlzdGljRGF0YUV4cG9ydCJdLCJpc0FkbWluIjpmYWxzZSwiaXNzIjoiZ2FzLXN0YXRpb24iLCJuYW1lIjoiY3VzdG9tZXIiLCJyb2xlcyI6W3sibmFtZSI6IuS6i-S7tuWkhOeQhuS6uuWRmCIsInJvbGVJRCI6IjE1ODc5MDg3NjYxNzc3MTQ3MjgifV0sInNpdGVJRHMiOltdfQ.84Re_Ss5Az11NQ_YaL9jn7dHbqSsKwTCfKbQDbkN588" // 请将其替换为实际的接口URL
	outputFile = "output.xlsx"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            // 输出的拼接后的Excel文件名
)

func main() {
	// 设置时间范围，这里以每小时为例，可以根据需求调整
	startTime := time.Date(2023, 7, 24, 0, 0, 0, 0, time.FixedZone("CST", 8*60*60))
	endTime := time.Date(2023, 7, 30, 23, 59, 59, 0, time.FixedZone("CST", 8*60*60))

	// 循环调用接口并拼接Excel文件
	err := fetchAndConcatenateExcel(startTime, endTime)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func fetchAndConcatenateExcel(startTime, endTime time.Time) error {
	f := excelize.NewFile()

	// 循环调用接口并拼接Excel文件
	for currentTime := startTime; currentTime.Before(endTime); currentTime = currentTime.Add(time.Hour) {
		d := DateTimeOrRange{From: sqlxDatatypes.Timestamp(currentTime), To: sqlxDatatypes.Timestamp(currentTime.Add(time.Hour))}
		ds, _ := d.MarshalText()
		dss := url.QueryEscape(string(ds))

		url := fmt.Sprintf("%s&passTimeRange=%s", apiURL, dss)

		fmt.Println(url)

		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
		}

		// 将返回的Excel文件内容拼接到输出文件
		if err := appendExcelToOutput(resp.Body, f); err != nil {
			return err
		}
	}

	// 保存拼接后的Excel文件
	if err := f.SaveAs(outputFile); err != nil {
		return err
	}

	fmt.Printf("拼接完成，文件保存为：%s\n", outputFile)
	return nil
}

func appendExcelToOutput(body io.Reader, f *excelize.File) error {
	xlsx, err := excelize.OpenReader(body)
	if err != nil {
		return err
	}

	// 复制每个Sheet到输出文件
	for _, sheetName := range xlsx.GetSheetMap() {
		rows, err := xlsx.GetRows(sheetName)
		if err != nil {
			return err
		}

		// 如果Sheet不存在，则创建新Sheet
		index, err := f.GetSheetIndex(sheetName)
		if err != nil {
			return err
		}
		if index == -1 {
			index, err = f.NewSheet(sheetName)
			if err != nil {
				return err
			}
		}

		f.SetActiveSheet(index)

		// 获取当前Sheet最后一行的索引
		// 获取当前Sheet最后一行的索引
		rowIndex := getRowsCount(f, sheetName) + 1

		// 追加新的内容到现有Sheet上
		for _, row := range rows {
			rowIndex++
			for colIndex, cellValue := range row {
				cellName, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex)
				if err != nil {
					return err
				}

				if err := f.SetCellValue(sheetName, cellName, cellValue); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// 获取指定Sheet的行数
func getRowsCount(f *excelize.File, sheetName string) int {
	rows, _ := f.GetRows(sheetName)
	if rows == nil {
		return 0
	}
	return len(rows)
}

type DateTimeOrRange struct {
	From sqlxDatatypes.Timestamp
	To   sqlxDatatypes.Timestamp
	ValueOrRangeOpt
}

func (timeRange *DateTimeOrRange) IsZero() bool {
	return timeRange == nil || (timeRange.From.IsZero() && timeRange.To.IsZero())
}

func (timeRange *DateTimeOrRange) ConditionFor(c *builder.Column) builder.SqlCondition {
	return timeRange.ValueOrRangeOpt.ConditionFor(c, &timeRange.From, &timeRange.To)
}

func (timeRange DateTimeOrRange) MarshalText() ([]byte, error) {
	return timeRange.ValueOrRangeOpt.MarshalText(&timeRange.From, &timeRange.To)
}

func (timeRange *DateTimeOrRange) UnmarshalText(data []byte) error {
	tr := DateTimeOrRange{}

	err := tr.ValueOrRangeOpt.UnmarshalText(data, &tr.From, &tr.To)
	if err != nil {
		return err
	}

	if !tr.From.IsZero() && !tr.To.IsZero() {
		if time.Time(tr.From).After(time.Time(tr.To)) {
			return fmt.Errorf("time from should not after time to")
		}
	}

	*timeRange = tr
	return nil
}

type Textable interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	IsZero() bool
}

// inspire by https://www.logicbig.com/tutorials/misc/groovy/range-operator.html
type ValueOrRangeOpt struct {
	ExclusiveFrom bool
	ExclusiveTo   bool
	Exactly       bool
}

func (valueOrRange *ValueOrRangeOpt) ConditionFor(c *builder.Column, from Textable, to Textable) builder.SqlCondition {
	where := builder.EmptyCond()

	if valueOrRange != nil {
		if !from.IsZero() {
			if valueOrRange.Exactly {
				return c.Eq(from)
			}

			if valueOrRange.ExclusiveFrom {
				where = where.And(c.Gt(from))
			} else {
				where = where.And(c.Gte(from))
			}
		}

		if !to.IsZero() {
			if valueOrRange.ExclusiveTo {
				where = where.And(c.Lt(to))
			} else {
				where = where.And(c.Lte(to))
			}
		}
	}

	return where
}

func (valueOrRange *ValueOrRangeOpt) UnmarshalText(text []byte, fromValue Textable, toValue Textable) error {
	if len(text) == 0 {
		return nil
	}

	r := ValueOrRangeOpt{}

	spliter := []byte("..")

	r.Exactly = !bytes.Contains(text, spliter)

	fromTo := bytes.Split(text, spliter)

	if len(fromTo) > 0 {
		from := fromTo[0]

		if len(from) > 0 {
			lastChar := from[len(from)-1]
			if lastChar == '!' || lastChar == '<' {
				from = from[:len(from)-1]
				r.ExclusiveFrom = true
			}
		}

		if len(from) > 0 {
			err := fromValue.UnmarshalText(from)
			if err != nil {
				return err
			}
		}
	}

	if len(fromTo) > 1 {
		to := fromTo[1]

		if len(to) > 0 {
			firstChar := to[0]
			if firstChar == '!' || firstChar == '<' {
				to = to[1:]
				r.ExclusiveTo = true
			}
		}

		if len(to) > 0 {
			err := toValue.UnmarshalText(to)
			if err != nil {
				return err
			}
		}
	}

	*valueOrRange = r

	return nil
}

func (valueOrRange ValueOrRangeOpt) MarshalText(fromValue Textable, toValue Textable) (text []byte, err error) {
	buf := bytes.NewBuffer(nil)

	if !fromValue.IsZero() {
		from, err := fromValue.MarshalText()
		if err != nil {
			return nil, err
		}

		buf.Write(from)
		if valueOrRange.ExclusiveFrom {
			buf.WriteByte('<')
		}
	}

	if !valueOrRange.Exactly {
		buf.WriteString("..")

		if !toValue.IsZero() {
			if valueOrRange.ExclusiveTo {
				buf.WriteByte('<')
			}
			to, err := toValue.MarshalText()
			if err != nil {
				return nil, err
			}

			buf.Write(to)
		}
	}

	return buf.Bytes(), nil
}
