package main

import (
	"fmt"
	"time"

	. "dedicate_server/gcore"
)

func main() {
	GetWorld()
	serverAddr := RunUDPServer("0.0.0.0", 51081)
	fmt.Println(serverAddr)
	for {
		time.Sleep(3 * time.Second)
	}
}
