package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// unsafe 包
	//确定结构在内存中占用的确切大小
	//func Sizeof(x ArbitraryType) uintptr
	//返回结构体中某个field的偏移量
	//func Offsetof(x ArbitraryType) uintptr
	//返回结构体中某个field的对其值（字节对齐的原因）
	//func Alignof(x ArbitraryType) uintptr

	// uintptr是一个无符号的整型数，足以保存一个地址，可以和unsafe.Pointer 相互转换
	// unsafe.Pointer 是一个指向变量的指针，因此当变量被移动时对应的指针也必须被更新
	// uintptr 类型的变量只是一个普通的数字

	a := [4]int{1, 2, 3, 4}
	ptr := (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(&a[1])) + unsafe.Sizeof(a[0])))
	*ptr = 9
	fmt.Println(a)

	h := struct {
		sex  bool
		age  uint8
		min  int
		name string
	}{
		true,
		30,
		1,
		"hello",
	}
	i := unsafe.Sizeof(h)
	j := unsafe.Alignof(h.age)
	k := unsafe.Offsetof(h.name)
	fmt.Println(i, j, k)
	fmt.Printf("%p\n", &h)
	var p unsafe.Pointer
	p = unsafe.Pointer(&h)
	fmt.Println(p)
}
