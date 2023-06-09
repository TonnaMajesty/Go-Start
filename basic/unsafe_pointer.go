package main

import (
	"fmt"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
)

type user struct {
	name string
	age  int
}

func main() {
	// case1
	// &l是一个指向l变量的指针，&l的值表示的是l变量的地址
	// unsafe.Pointer(&l)表示的是l变量的地址。unsafe.Pointer(&l)是将这个指针转换为一个unsafe.Pointer类型的值。
	l := 5
	fmt.Println(unsafe.Pointer(&l)) // 0x1400012c160
	spew.Dump(&l)                   // (*int)(0x1400012c160)(5)

	// case 2
	// 创建一个user对象
	u := user{"Alice", 20}
	fmt.Println(u) // {Alice 20}

	// 获取user对象的地址，并转换为unsafe.Pointer
	p := unsafe.Pointer(&u)
	fmt.Println(p) // 0xc00000a0a0

	// 获取user对象中name字段的地址，并转换为unsafe.Pointer
	pname := unsafe.Pointer(uintptr(p) + unsafe.Offsetof(u.name))
	fmt.Println(pname) // 0xc00000a0a0

	// 将pname转换为*string类型，并修改其值
	*(*string)(pname) = "Bob"
	fmt.Println(u) // {Bob 20}

	// 获取user对象中age字段的地址，并转换为unsafe.Pointer
	page := unsafe.Pointer(uintptr(p) + unsafe.Offsetof(u.age))
	fmt.Println(page) // 0xc00000a0b8

	// 将page转换为*int类型，并修改其值
	*(*int)(page) = 30
	fmt.Println(u) // {Bob 30}
}
