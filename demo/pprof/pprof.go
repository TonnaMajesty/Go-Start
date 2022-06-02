package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)
var datas []string

func Add(str string) string {
	data := []byte(str)
	sData := string(data)
	datas = append(datas, sData)

	return sData
}


func main() {
	var mu sync.Mutex
	count := 100
	go func() {
		for i:=1;i<10000;i++ {
			mu.Lock()
			count++
			mu.Unlock()
		}
	}()

	go func() {
		for i:=1;i<10000;i++ {
			mu.Lock()
			count--
			mu.Unlock()
		}
	}()

	go func() {
		for {
			log.Println(Add("https://github.com/EDDYCJY"))
			time.Sleep(time.Second)
		}
	}()

	http.ListenAndServe("0.0.0.0:6060", nil)
}
