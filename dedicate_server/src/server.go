package main

import (
	"fmt"
	"time"

	"dedicate_server/gcore"
	"dedicate_server/gnet"
)

func main() {
	gcore.GetWorld()
	serverAddr := gnet.RunUDPServer("127.0.0.1", 51081)
	fmt.Println(serverAddr)
	for {
		time.Sleep(3 * time.Second)
	}
}
