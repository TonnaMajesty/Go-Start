package main

import (
	"sync"
)

type test struct {
	mutex sync.Mutex
}

func main() {
	t := test{}
	t.mutex.Lock()
}
