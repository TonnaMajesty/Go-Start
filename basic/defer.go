package main

import "fmt"

func test() (i int) { // 命名返回值，函数内部的所有修改都会影响返回值
	i = 0
	defer func() {
		fmt.Println("defer1")
	}()
	defer func() {
		i += 1
		fmt.Println("defer2")
	}()
	return i
}

func main() {
	fmt.Println("return", test())
}