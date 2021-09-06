package main

import (
	"fmt"
	"time"

	"dedicate_server/gnet"
)

func main() {
	serverAddr := gnet.MakeUDPServer("127.0.0.1", 50080)
	fmt.Println(serverAddr)
	for {
		time.Sleep(3 * time.Second)
	}
}
