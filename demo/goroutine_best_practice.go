package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	tr := NewTracker()
	go tr.Run()
	_ = tr.Event(context.Background(), "test1")
	_ = tr.Event(context.Background(), "test2")
	_ = tr.Event(context.Background(), "test3")
	//time.Sleep(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	tr.Shutdown(ctx)
}

func NewTracker() *Tracker {
	return &Tracker {
		ch: make(chan string, 10),
	}
}

type Tracker struct {
	ch chan string
	stop chan struct{}
}

func (t *Tracker) Event (ctx context.Context, data string) error {
	select {
		case t.ch <- data:
			return nil
		case <-ctx.Done():
			return ctx.Err()
	}
}

func (t *Tracker) Run () {
	for data := range t.ch { // range 从channel中取数据，如果channel为空也会阻塞，channel关闭循环自动退出
		//time.Sleep(time.Second)
		fmt.Println(data)
	}

	t.stop <- struct{}{}
}

func (t *Tracker) Shutdown(ctx context.Context) {
	close(t.ch)
	select {
	case <-t.stop:
	case <-ctx.Done():
	}
}