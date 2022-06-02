package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx1, cancel1 := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	ctx2, _ := context.WithTimeout(ctx1, 5*time.Second)
	ctx2 = context.WithValue(ctx2, "hello", "world")



	go func(ctx context.Context) {
		defer func() {
			fmt.Println("exit1")
		}()
		select {
			case <-ctx.Done(): // withDeadline 到期或者手动调用了cancel，Done都会收到通知
				fmt.Println("done1")
			case <-time.After(10*time.Second):
				fmt.Println("timeout1")
		}
	}(ctx1)

	go func(ctx context.Context) {
		defer func() {
			fmt.Println("exit2")
		}()
		select {
		case <-ctx.Done(): // ctx2 可以接收到 ctx1的done通知
			fmt.Println("done1")
			fmt.Println(ctx.Err())
		case <-time.After(10*time.Second):
			fmt.Println("timeout2")
		}
	}(ctx2)

	d, _ := ctx1.Deadline()
	fmt.Println(d)
	v := ctx2.Value("hello")
	fmt.Println(v)

	time.Sleep(time.Second)
	cancel1()
	time.Sleep(5*time.Second)
}
