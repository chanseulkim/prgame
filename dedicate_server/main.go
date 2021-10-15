package main

import (
	"fmt"
	"time"

	. "dedicate_server/gcore"
)

func main() {
	GetWorld()
	serverAddr := RunTcpServer("0.0.0.0", 51080)
	fmt.Println(serverAddr)
	for {
		time.Sleep(3 * time.Second)
	}
}
