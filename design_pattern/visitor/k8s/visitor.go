package k8s

import (
	"fmt"
)

// 装饰器模式 + 访问者模式

// 访问者
type VisitorFunc func(*Info, error) error

// Info 是被访问者
type Info struct {
	Namespace string
	Name string
	OtherThings string
}
// Info接收一个访问者方法
func (info *Info) Accept (fn VisitorFunc) error {
	return fn(info, nil)
}

// 被访问者装饰器, 由于被访问者最终要传入装饰器，所以被访问者和装饰器使用同一个interface
type DecoratedVisitor interface {
	Accept(VisitorFunc) error
}

// 这个Visitor 主要是用来访问 Info 结构中的 Name 和 NameSpace 成员
type NameVisitor struct {
	decoratedVisitor DecoratedVisitor
}

func (v NameVisitor) Accept(fn VisitorFunc) error {
	return v.decoratedVisitor.Accept(func(info *Info, err error) error {
		fmt.Println("NameVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> Name=%s, NameSpace=%s\n", info.Name, info.Namespace)
		}
		fmt.Println("NameVisitor() after call function")
		return err
	})
}

// 这个Visitor主要用来访问 Info 结构中的 OtherThings 成员
type OtherThingsVisitor struct {
	decoratedVisitor DecoratedVisitor
}

func (v OtherThingsVisitor) Accept(fn VisitorFunc) error {
	return v.decoratedVisitor.Accept(func(info *Info, err error) error {
		fmt.Println("OtherThingsVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("==> OtherThings=%s\n", info.OtherThings)
		}
		fmt.Println("OtherThingsVisitor() after call function")
		return err
	})
}

// log visitor
type LogVisitor struct {
	decoratedVisitor DecoratedVisitor
}

func (v LogVisitor) Accept(fn VisitorFunc) error {
	return v.decoratedVisitor.Accept(func(info *Info, err error) error {
		fmt.Println("LogVisitor() before call function")
		err = fn(info, err)
		fmt.Println("LogVisitor() after call function")
		return err
	})
}





