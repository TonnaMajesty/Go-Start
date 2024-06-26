package main

type ifoo interface {
	bar()
}

type foo struct {
}

func (f foo) bar() {

}

type fooo struct {
}

func (f *fooo) bar() {

}

// 结构体实现了接口，结构体指针就实现了接口
var _ ifoo = &foo{}

// 指针实现了接口不代表结构体实现了接口
//var _ ifoo = fooo{} // 错误

func main() {
	// 结构体和指针间的方法可以相互调用
	foo1 := &foo{}
	foo1.bar()

	foo2 := foo{}
	foo2.bar()
}
