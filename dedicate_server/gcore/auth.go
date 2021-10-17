package gcore

import (
	"fmt"
	. "gnet"
	"net"
)

func handleEnterClient(user_id string, client_addr net.Addr, pos Vector2) bool {
	_, exists := GetWorld().Players[user_id]
	if exists == false {
		fmt.Println("enter client : " + client_addr.String() + ", " + user_id)
		GetWorld().AddPlayer(user_id, client_addr, pos)
		GetWorld().Players[user_id] = NewPlayer(0, user_id, client_addr, pos, DEFAULT_COLISION_RADIUS)
	}
	return true
}
