package main

import (
	"fmt"
	"time"

	"github.com/juju/ratelimit"
)

func main() {
	for i := 0; i < 100; i++ {
		Do(i)
		time.Sleep(time.Millisecond * 400)
	}
}

// 2 tokens per second, and a maximum burst of 2.
var bucket = ratelimit.NewBucketWithRate(2, 2)

func Do(i int) {
	ok := bucket.WaitMaxDuration(1, 0)
	if !ok {
		fmt.Printf("Request %d failed to get token\n", i)
		return
	}

	fmt.Printf("Request %d took success \n", i)
}
