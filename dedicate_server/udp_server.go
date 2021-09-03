package main

import (
	"fmt"
	"time"

	"./core_libs"
)

func main() {
	serverAddr := core_libs.MakeUDPServer("127.0.0.1", 50080)
	fmt.Println(serverAddr)
	for {
		time.Sleep(3 * time.Second)
	}

}
