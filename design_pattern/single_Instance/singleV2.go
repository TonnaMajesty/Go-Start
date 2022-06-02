package main

import (
	"fmt"
	"sync"
)

// 另外两个思路
// 1. 在init函数中创建实例
// 2. 使用sync.Once

var once sync.Once

type singleV2 struct {
}

var singleInstanceV2 *singleV2

func getInstanceV2() *singleV2 {
	if singleInstanceV2 == nil {
		once.Do(
			func() {
				fmt.Println("Creating single instance now.")
				singleInstanceV2 = &singleV2{}
			})
	} else {
		fmt.Println("Single instance already created.")
	}

	return singleInstanceV2
}
