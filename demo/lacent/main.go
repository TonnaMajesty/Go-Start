package main

import (
	"fmt"

	datastructure "github.com/duke-git/lancet/v2/datastructure/queue"
)

func main() {

	q := datastructure.NewArrayQueue[int](5)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	fmt.Println(q.Dequeue())
}
