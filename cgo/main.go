package main

// 在C语言中，修饰符extern用在变量或者函数的声明前，用来说明“此变量/函数是在别处定义的，要在此处引用”

//extern int go_qsort_compare(void* a, void* b);
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

// https://chai2010.cn/advanced-go-programming-book/ch2-cgo/ch2-06-qsort.html

//export go_qsort_compare
func go_qsort_compare(a, b unsafe.Pointer) C.int {
	pa, pb := (*C.int)(a), (*C.int)(b)
	return C.int(*pa - *pb)
}

// 闭包函数无法导出为 C 语言函数，因此无法直接将闭包函数传入 C 语言的 qsort 函数。
// 为此我们可以用 Go 构造一个可以导出为 C 语言的代理函数，然后通过一个全局变量临时保存当前的闭包比较函数。
var go_qsort_compare_info struct {
	// SortV2 相关
	sync.Mutex
	fn func(a, b unsafe.Pointer) int

	// Slice 相关
	base     unsafe.Pointer
	elemnum  int
	elemsize int
	less     func(a, b int) bool
}

//export _cgo_qsort_compare
func _cgo_qsort_compare(a, b unsafe.Pointer) C.int {
	return C.int(go_qsort_compare_info.fn(a, b))
}

func main() {
	values := []int32{42, 9, 101, 95, 27, 25}

	Sort(unsafe.Pointer(&values[0]),
		len(values), int(unsafe.Sizeof(values[0])),
		CompareFunc(C.go_qsort_compare),
	)
	fmt.Println(values)

	// 改进：闭包函数作为比较函数
	SortV2(unsafe.Pointer(&values[0]), len(values), int(unsafe.Sizeof(values[0])),
		func(a, b unsafe.Pointer) int {
			pa, pb := (*int32)(a), (*int32)(b)
			return int(*pa - *pb)
		},
	)

	fmt.Println(values)

	// 继续改进包装函数，尝试消除对 unsafe 包的依赖, 实现一个类似标准库中 sort.Slice 的排序函数
	Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	fmt.Println(values)
}
