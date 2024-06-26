package main

import (
	"fmt"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
)

type person struct {
	Name string
	Aget string
}

func main() {
	// case1
	// 映射、切片共享底层元素，修改会互相影响
	m := map[string]string{"name": "tzh"}
	s := make([]map[string]string, 0)

	s = append(s, m)
	fmt.Println(s)
	m["name"] = "mqb"
	fmt.Println(s)

	// case2
	// 这段代码中pMap的结果是相同的，是因为在循环中使用了指针&l来赋值给map的value，
	// l是一个临时变量，每次循环值都会被覆盖，但是l的地址在每次循环中都是不变的，所以map的value都会指向最后一次循环时l的地址和值。
	list := []person{{"tzh", "20"}, {"mqb", "22"}}
	pMap := map[string]*person{}
	for _, l := range list {
		fmt.Println(unsafe.Pointer(&l))
		// 以在循环内部重新声明一个变量来存储l的值，然后用指针赋值给map的value，这样就可以避免这个问题。
		//l := l
		pMap[l.Name] = &l
	}

	spew.Println(pMap)

	// case3 修改map中结构体的值
	type Person struct {
		Name string
		Age  int
	}

	people := make(map[string]Person)
	people["Alice"] = Person{"Alice", 25}
	people["Bob"] = Person{"Bob", 30}

	// 不可以。这样做会导致编译错误，因为 map 中的值是不可寻址的(当通过key获取到value时，这个value是不可寻址的)。
	// 也就是说，你不能直接修改 map 中的结构体
	// people["Alice"].Age = 26

	// https://forum.golangbridge.org/t/question-about-modifying-struct-in-map/23557/2
	// It’s because map keys or values are not addressable
	// If you add or remove an element to a map, the elements can be shuffled around in the backing memory.
	// If you could take the address of an element and then modified the map, who knows what your pointer would point to which value?
	// maybe the element associated with some other key, maybe some sentinel value to indicate that the slot is available, etc.

	// maps of pointers to structs have the same issue, but in that case
	// it’s just the pointers that get shuffled around when values are added and not the whole embedded structs.

	people2 := make(map[string]*Person)
	people2["Alice"] = &Person{"Alice", 25}
	people2["Bob"] = &Person{"Bob", 30}
	// 这里map的值存储的是Person的指针，即使map因为扩容被打乱，指针指向的Person是不会受影响的，所以可以通过指针对Person的值进行修改
	people2["Alice"].Age = 26
}
