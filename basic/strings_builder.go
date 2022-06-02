package main

import (
	"fmt"
	"strings"
)

func main() {
	var s strings.Builder
	for i := 0;i<10;i++ {
		s.WriteString("hello")
	}

	fmt.Println(s.String())
}
