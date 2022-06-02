package main

import (
	"errors"
	"fmt"
	"os"
)


func main() {
	var err error
	err = &os.PathError{
		Path: "err path",
		Op:   "err op",
		Err:  errors.New("path err"),
	}
	// 无法确定wrappedErr的类型
	wrappedErr := fmt.Errorf("wrappedErr %w", err)
	if _, ok := wrappedErr.(*os.PathError); ok {
		fmt.Println("assert wrappedErr  os.PathError")
	}
	// 用As来断言wrappedErr的类型
	var perr *os.PathError
	if errors.As(wrappedErr, &perr) {
		fmt.Println("as wrappedErr os.PathError")
	}


	// error is 递归调用 Unwrap 并判断每一层的 err 是否相等，如果有任何一层 err 和传入的目标错误相等，则返回 true
	err1 := errors.New("new error")
	err2 := fmt.Errorf("err2: [%w]", err1)
	err3 := fmt.Errorf("err3: [%w]", err2)

	fmt.Println(errors.Is(err3, err2))
	fmt.Println(errors.Is(err3, err1))


	// error as 和上面的 errors.Is 大体上是一样的，区别在于 Is 是严格判断相等，即两个 error 是否相等。
	// 而 As 则是判断类型是否相同，并提取第一个符合目标类型的错误，用来统一处理某一类错误
	var targetErr *ErrorString
	err4 := fmt.Errorf("new error:[%w]", &ErrorString{s:"target err"})
	fmt.Println(errors.As(err4, &targetErr))

}

type ErrorString struct {
	s string
}

func (e *ErrorString) Error() string{
	return e.s
}


