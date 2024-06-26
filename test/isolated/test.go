package isolated

import (
	"fmt"
)

type Inner2 struct {
}

func (i Inner2) Init() {
	fmt.Println("inner2 init")

}
