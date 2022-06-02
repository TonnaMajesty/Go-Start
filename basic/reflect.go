package main

import (
	"fmt"
	"reflect"
)

type order struct {
	Id int `json:"id"`
	Orderid int `gorm:"id"`
	Userid int
	Username string
}

func main() {
	o := order{123, 234, 678, "tonna"}
	v := reflect.ValueOf(o)
	t := reflect.TypeOf(o)
	k := t.Kind()
	k1 := v.Kind()
	fmt.Printf("%#v\n", v)
	fmt.Printf("%#v\n", t)
	fmt.Printf("%v\n", k)
	fmt.Printf("%v\n", k1)
	for i:=0; i < t.NumField(); i++ {
		fmt.Println("FieldName:", t.Field(i).Name, "FiledType:", t.Field(i).Type,
			"FiledValue:", v.Field(i), "FiledTag:", t.Field(i).Tag)
	}
}
