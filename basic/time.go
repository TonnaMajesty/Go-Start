package main

import (
	"fmt"
	"time"
)

func main() {
	// 获取当前时间
	fmt.Println(time.Now())
	fmt.Println(time.Now().Year())
	fmt.Println(time.Now().Day())

	// Format
	fmt.Println(time.Now().Format("2006-01-02 15:04:05")) // 2006 年 1 月 2 日 3 点 4 分 5

	// 设置时区
	location, _ := time.LoadLocation("Asia/Shanghai")
	fmt.Println(time.Now().In(location))

	// 获取时间戳
	fmt.Println(time.Now().Unix())
	fmt.Println(time.Now().UnixNano()) // 纳秒时间戳

	// 时间戳转换为时间
	time.Unix(time.Now().Unix(), 0) //将时间戳转为时间格式

	//在字符串中解析出 duration（持续时间）
	//其支持的有效单位有"ns”, “us” (or “µ s”), “ms”, “s”, “m”, “h”，例如：“300ms”, “-1.5h” or “2h45m”
	duration, err := time.ParseDuration("100m")
	if err != nil {
		fmt.Println(duration)
	}

	// Parse
	inputTime := "2029-09-04 12:02:33"
	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, inputTime)
	fmt.Println(t)
	tl, _ := time.ParseInLocation(layout, inputTime, location)
	fmt.Println(tl)

	// Add
	fmt.Println(time.Now().Add(duration))
	fmt.Println(time.Now().Add(time.Second * 20))

	// sub
	fmt.Println(time.Now().Sub(time.Now().Add(time.Second)))

	// before after 返回true false
	now := time.Now()
	later := now.Add(time.Minute)

	fmt.Println(later.Before(now))
	fmt.Println(later.After(now))

	// 定时器 Tick 本质上是一个channel
	ticker := time.Tick(time.Second) //定义一个1秒间隔的定时器, 返回 <- chan Time
	for i := range ticker {
		fmt.Println(i) //每秒都会执行的任务
	}

	ticker2 := time.NewTicker(time.Second)
	defer ticker2.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second)
		done <- true
	}()
	for {
		select {
		case <-done:
			fmt.Println("Done!")
		case t := <-ticker2.C:
			fmt.Println("Current time: ", t)     }
	}

}
