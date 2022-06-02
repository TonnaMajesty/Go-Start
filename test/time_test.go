package test

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("bigin")
	m.Run()
	fmt.Println("end")
}

func TestTime(t *testing.T) {
	t.Error("testtime")
}

func TestHello(t *testing.T) {
	t.Error("hello")
}
