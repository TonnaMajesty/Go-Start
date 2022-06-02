package main

import "fmt"

type Parent struct {

}

func (p Parent) SayHello() {
	fmt.Println("I am " + p.Name())
}

func (p Parent) Name() string {
	return "Parent"
}

type Son struct {
	Parent
}

func (s Son) Name() string {
	return "Son"
}

func main() {
	son := Son{
		Parent{},
	}

	son.SayHello()
}
