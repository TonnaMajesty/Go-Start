package main

/*
#include <errno.h>
// C 语言不支持返回多个结果，因此 <errno.h> 标准库提供了一个 errno 宏用于返回错误状态。
// 我们可以近似地将 errno 看成一个线程安全的全局变量，可以用于记录最近一次错误的状态码。
static int div(int a, int b) {
    if(b == 0) {
        errno = EINVAL;
        return 0;
    }
    return a/b;
}
*/
import "C"
import (
	"fmt"
)

func DoDiv() {
	// CGO 也针对 <errno.h> 标准库的 errno 宏做的特殊支持：在 CGO 调用 C 函数时如果有两个返回值，那么第二个返回值将对应 errno 错误状态。
	v0, err0 := C.div(2, 1)
	fmt.Println(v0, err0)

	v1, err1 := C.div(1, 0)
	fmt.Println(v1, err1)
}
