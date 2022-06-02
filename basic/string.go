package main

import (
	"fmt"
	"strings"
	"unicode"
)

func main() {
	// split
	s := "123.456.789.3434"
	s1 := strings.SplitN(s, ".", -1) // -1 全部
	fmt.Println(s1)
	s2 := strings.Split(s, ".")
	fmt.Println(s2)

	// join
	sl := []string{"123", "456", "789"}
	s3 := strings.Join(sl, "-")
	fmt.Println(s3)

	// count
	s4 := "helloworld"
	fmt.Println(strings.Count(s4, "o")) // o出现的次数
	fmt.Println(strings.Count(s4, "")) // 字符串长度+1
	// Index
	fmt.Println(strings.Index(s4, "h"))
	// replace
	s4 = strings.Replace(s4, "o", "1", -1)
	fmt.Println(s4)

	// unicode.IsLetter 类似的还有unicode.IsNumber
	for _,v := range "h1ell2o" {
		fmt.Println(unicode.IsLetter(v)) // v 是 rune类型
	}
}
