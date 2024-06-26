package main

import (
	"fmt"
	"time"

	"github.com/tbrandon/mbserver"
)

func main() {
	serv := mbserver.NewServer()
	err := serv.ListenTCP("127.0.0.1:1502")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	defer serv.Close()

	// Wait forever
	for {
		time.Sleep(1 * time.Second)
	}
}
