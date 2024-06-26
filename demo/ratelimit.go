package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter interface {
	Wait(context.Context) error
	Limit() rate.Limit
}

type multiLimiter struct {
	limiters []RateLimiter
}

func MultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}
func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}
func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}

func main() {
	//1秒钟2个
	secondLimit := rate.NewLimiter(Per(5, 1*time.Second), 1)
	//60秒100个
	minuteLimit := rate.NewLimiter(Per(100, 1*time.Minute), 20)
	multiLimiter := MultiLimiter(secondLimit, minuteLimit)

	for {
		{
			if err := multiLimiter.Wait(context.Background()); err == nil {
				time.Sleep(1 * time.Second)
				fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			} else {
				fmt.Println("err")
			}
		}
	}

}
func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}
