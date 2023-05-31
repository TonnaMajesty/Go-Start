package main

import (
	"encoding/json"
	"fmt"

	"github.com/coocood/freecache"
)

type a struct {
	Name string `json:"name,omitempty"`
}

func main() {
	// In bytes, where 1024 * 1024 represents a single Megabyte, and 100 * 1024*1024 represents 100 Megabytes.
	cacheSize := 100 * 1024 * 1024
	cache := freecache.NewCache(cacheSize)

	key := []byte("abc")
	val, _ := json.Marshal(a{"tzh"})

	expire := 60 // expire in 60 seconds
	cache.Set(key, val, expire)
	got, err := cache.Get(key)

	a := &a{}
	json.Unmarshal(got, a)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v", a)
	}
	affected := cache.Del(key)
	fmt.Println("deleted key ", affected)
	fmt.Println("entry count ", cache.EntryCount())
}
