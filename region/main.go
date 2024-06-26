package main

import (
	"fmt"
)

func main() {
	// 打印查看
	region := NewGBT2260()
	localCode := region.SearchGBT2260("130104")
	fmt.Println(localCode)

	province := region.GetAllProvince()
	fmt.Println(province)
}
