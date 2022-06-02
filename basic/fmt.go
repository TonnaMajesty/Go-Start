package main

import "fmt"

type Stu struct {
	Name string
}


// 但如果结构体定义了 String() 方法，%v 和 %+v 都会调用 String() 覆盖默认值
func (s Stu) String() string{
	return "123"
}

func main() {
	// %v 和 %+v 都可以用来打印 struct 的值，区别在于 %v 仅打印各个字段的值，%+v 还会打印各个字段的名称
	fmt.Printf("%v\n", Stu{"Tom"}) // {Tom}
	fmt.Printf("%+v\n", Stu{"Tom"}) // {Name:Tom}
}
