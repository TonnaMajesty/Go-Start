package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	eg, _ := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		time.Sleep(5 * time.Second)
		err := errors.New("go1 err")
		return err
		//select {
		//case <-ctx.Done():
		//	fmt.Println("go1 cancel, err = ", ctx.Err())
		//default:
		//	fmt.Println("go1 run")
		//}
		//return nil
	})
	eg.Go(func() error {
		err := errors.New("go2 err")
		return err
	})
	err := eg.Wait()
	if err != nil {
		fmt.Println("err =", err)
	}
}
