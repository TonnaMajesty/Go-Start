package main

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

func main() {
	m := map[string]interface{}{
		"names": "123",
	}

	d := P{}

	mapstructure.Decode(m, &d)
	fmt.Println(d)

}

type P struct {
	NameS string `mapstructure:"names"`
}
