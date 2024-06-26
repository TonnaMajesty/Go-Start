package main

import (
	"fmt"
	"time"
)

func main() {
	region, _ := NewRegion()

	province := region.GetAllProvince()
	fmt.Println(province)

	fmt.Println(region.GetCityByProvince("51"))

	fmt.Println(region.SearchCodeByName("成都市"))

	cell, _ := region.SearchByCode("511902")
	fmt.Println(cell)

	fmt.Println(region.GetDistrictByCity("5101"))

	now := time.Now()
	res := []Cell{}
	for _, code := range DivideCode("632626") {
		c, _ := region.SearchByCode(code)
		res = append(res, c)
	}
	fmt.Println(res)
	fmt.Println(time.Now().UnixMilli() - now.UnixMilli())

	//spew.Println(region.trie.root.children["11"].children["01"])
}
