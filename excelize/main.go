package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("/Users/tonnamajesty/Downloads/充电站监管信息表 (社会桩).xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// 关闭文件
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetList := f.GetSheetList()
	//// 从指定的单元格中取值
	//cell, err := f.GetCellValue(sheetList[0], "充电站编号")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(cell)
	// 从sheet中获取行数据
	rows, err := f.GetRows(sheetList[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}
