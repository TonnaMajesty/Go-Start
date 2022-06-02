package main

import (
	"errors"
	"fmt"
	xerrors "github.com/pkg/errors"
)

type MyError struct {
	s string
}

func (m *MyError) Error() string {
	return m.s
}

func main() {
	myerr := &MyError{"myerr"}
	myerr2 := fmt.Errorf("%w 错误", myerr)
	myerr3 := xerrors.Wrap(myerr2, "myerr3")

	//var mye *MyError
	fmt.Println(xerrors.Is(myerr3, myerr))



	err1 := errors.New("错误")
	err2 := xerrors.Wrap(err1, "message")

	if xerrors.Is(err2, err1) {
		fmt.Println("mmmmmmm")
	}

	fmt.Printf("%T", err2)



	//err := xerrors.New("hello")
	//err = xerrors.Errorf("%w world", err)
	//fmt.Println(xerrors.Cause(err))
	//spew.Dump(err)
}