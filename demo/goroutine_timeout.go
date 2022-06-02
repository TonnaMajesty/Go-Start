package main

import (
	"context"
	"errors"
	"time"
)

// 模拟一个耗时操作
func search(term string) (string, error) {
	time.Sleep(200*time.Millisecond)
	return "some value", nil
}

//func process(term string) error {
//	record, err := search(term)
//	if err != nil {
//		return err
//	}
//
//	fmt.Println("Received", record)
//	return nil
//}

type result struct {
	record string
	err error
}

func process(term string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ch := make(chan result)

	go func() {
		record, err := search(term)
		ch <-result{record: record, err: err}
	}()

	select {
	case <-ctx.Done():
		return errors.New("search canceled")
	case result := <-ch:
		if result.err != nil {
			return result.err
		}
		return nil
	}
}


