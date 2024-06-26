package main

import (
	"fmt"
)

func main() {
	n := 10
	cameraNum := 5
	for i := 1; i <= n; i++ {
		startBottle := (i-1)*cameraNum + 1
		for i1 := 1; i1 <= cameraNum; i1++ {
			currentBottle := startBottle + (cameraNum - i1)
			fmt.Printf("相机：%d 瓶子：%d\n", i1, currentBottle)
		}
	}
}
