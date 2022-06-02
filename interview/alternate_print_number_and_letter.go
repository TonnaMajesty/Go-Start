package main

import (
	"fmt"
	"strings"
	"sync"
)

// 交替打印数字和字母  12AB34CD56EF78GH
func main() {
	letter, number := make(chan bool), make(chan bool)
	go func() {
		i := 1
		for {
			<-number
			fmt.Print(i)
			i++
			fmt.Print(i)
			i++
			letter <- true

		}
	}()

	var wait sync.WaitGroup
	wait.Add(1)

	go func(wait *sync.WaitGroup) {
		str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		i := 0
		for {
			<-letter
			if i >= strings.Count(str, "")-1 { // strings.Count 计算s中substr的非重叠实例的数量。如果substr是一个空字符串，Count返回1+s中Unicode码位的数量
				wait.Done()
				return
			}
			fmt.Print(str[i : i+1])
			i++
			fmt.Print(str[i : i+1])
			i++
			number <- true
		}

	}(&wait)
	number <- true
	wait.Wait()
}
