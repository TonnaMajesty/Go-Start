package main

import (
	"fmt"
	"sync/atomic"
)


type config struct {
	a []int
}
func main() {
	var v atomic.Value
	v.Store(config{a: []int{1,2,3}})
	cfg := v.Load().(config)
	fmt.Println(cfg)
}
