package main

import (
	"errors"
	"fmt"
	"time"
)

// 模拟了在goroutine的生命管理
// 在goroutine外部可以获取goroutine退出信号，也可以通知goroutine退出

func main() {
	stop := make(chan struct{})
	done := make(chan error)
	go func() {
		done <- server(stop)
	}()

	// 5秒后结束
	go func() {
		time.Sleep(5*time.Second)
		stop <- struct{}{}
	}()

	var stoped bool
	for i:=0; i<cap(done); i++ {
		if err := <-done; err != nil {
			fmt.Println(err)
		}
		if !stoped {
			stoped = true
			close(stop)
		}
	}
	time.Sleep(10*time.Second)
}

func server (stop chan struct{}) error{
	ch := make(chan struct{})
	go func() {
		<-stop
		close(ch)
	}()

	return listen(ch)
}

func listen(ch chan struct{}) error{
ForEnd:
	for {
		select {
		case <-ch:
			break ForEnd
		default:
		}
		time.Sleep(time.Second)
		fmt.Println(time.Now().UnixNano())

	}
	fmt.Println("listen exit")


	return errors.New("exit")
}

