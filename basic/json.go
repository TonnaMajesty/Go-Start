package main

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
)

type employee struct {
	Id int `json:"id,omitempty"`
	Name *string `json:"name,omitempty"` // 定义为指针类型，不会忽略零值，因为指针的零值是nil
}

func main() {
	st := new(string)
	*st = ""
	s1, err := json.Marshal(employee{Id:0, Name:st}) // marshal返回的是[]byte 类型
	s2, err := json.Marshal([]interface{}{employee{Id:0, Name:st}}) // 传入slice

	us1 := employee{}
	err = json.Unmarshal(s1, &us1) // unmarshal 传入的是byte类型
	us2 := make([]interface{}, 0)
	err = json.Unmarshal(s2, &us2)

	if err != nil {
		panic(err)
	}
	spew.Dump(string(s1))
	spew.Dump(string(s2))
	spew.Dump(us1)
	spew.Dump(us2)
}
