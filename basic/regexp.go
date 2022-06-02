package main

import (
	"fmt"
	"regexp"
)

func main() {
	re, _ := regexp.Compile(`a=(\d+),b=(\d+)`)
	s1 := re.ReplaceAllString("test regexp a=1234,b=5678. test regexp replace a=8765,b=3210 ", "c=$2,d=$1")
	fmt.Println(s1)
	s2 := re.FindAllString("test regexp a=1234,b=5678. test regexp replace a=8765,b=3210 ", -1)
	fmt.Println(s2)
	s3 := re.Split("test regexp a=1234,b=5678. test regexp replace a=8765,b=3210 ", -1)
	fmt.Println(s3)
}
