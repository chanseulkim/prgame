package main

import (
	"fmt"
	. "gnet"
)

type conninfo struct {
	user_id     string
	server_ip   string
	server_port string
}

func main() {
	JoinRoom("", "51080", "", "id", Vector2{X: 1, Y: 1})
	fmt.Println("gnetc")
}
