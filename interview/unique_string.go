package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "helowrdd"
	fmt.Println(isUniqueString(s))
	fmt.Println(isUniqueString2(s))

}

func isUniqueString(s string) bool {
	for k, v := range s {
		if strings.Index(s, string(v)) != k { // v 是 rune类型
			return false
		}
	}

	return true
}

func isUniqueString2(s string) bool {
	for _, v:= range s {
		if strings.Count(s, string(v)) > 1  {
			return false
		}
	}
	return true
}
