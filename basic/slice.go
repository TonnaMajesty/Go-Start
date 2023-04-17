package main

import (
	"fmt"
)

func RemoveFirstEle(s []any) any {
	copy(s, s[1:])
	fmt.Println(s)
	s = s[:len(s)-1]
	return s
}

func main() {
	fmt.Println(RemoveFirstEle([]any{"hello", "thank", "you"}))
}
