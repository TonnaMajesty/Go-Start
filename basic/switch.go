package main

import (
	"fmt"
	"math/rand"
)

func main() {
	x := 10
	switch x {
	case 1:
		fmt.Println("1") // 可以不写break
	case 10, 20, 30: // 包含多个表达式
		fmt.Println("10, 20, 30")
	default:
		fmt.Println("default")
	}

	switch {
	case x == 10:
		fmt.Println("10")
	default:
		fmt.Println("default")
	}

	// fallthrough
	//在 Go 中执行完一个 case 之后会立即退出 switch 语句。
	//fallthrough语句用于标明执行完当前 case 语句之后按顺序执行下一个case 语句。
	switch x := rand.Intn(100); { // 注意要加分号 类似if (v.(type) 不需要加分号)
	case x > 10:
		fmt.Println(">10")
		fallthrough
	case x > 20: // fallthrough 后的case语句不会判断是否满足，直接执行
		fmt.Println(">20")


	}

	var d animal = Dog{} // d:=Dog{} 这么写是不行的
	// 类型判断 type switch
	switch d.(type) { // d必须是接口类型的变量，所有case语句后面的类型必须实现了d的接口类型
	case Dog:
		fmt.Println(d.shout())
	}

}

type animal interface {
	shout() string
}

type Dog struct {

}

func (d Dog) shout() string {
	return "dog"
}
