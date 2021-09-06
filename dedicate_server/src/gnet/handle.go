package gnet

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	. "dedicate_server/gcore"
)

var sight_value Float = 20.0
var msg_queue_ch = make(chan *MsgBuff)

func handleJoin(userid string, client_addr net.Addr) {

	_, exists := GetWorld().Players[client_addr.String()]
	if exists == false {
		for _, player := range GetWorld().Players {
			usrid := player.Uid
			syncmsg := usrid + ";sync;" + last_position[usrid] + ";"
			syncmsg_len := len(syncmsg)
			syncmsg_buff := make([]byte, syncmsg_len)
			copy(syncmsg_buff, syncmsg)
			_, err := server.WriteTo(syncmsg_buff[:syncmsg_len], client_addr)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("enter client : " + client_addr.String() + ", " + userid)
		client_addrs[userid] = client_addr
		GetWorld().Players[client_addr.String()] = Player{Uid: userid}
	}
}

func handleMove(x Float, y Float, action string) {
	findFov(x, y, action)
}

func findFov(x Float, y Float, action string) {
	if action == "ui_left" {

	} else if action == "ui_right" {
	} else if action == "ui_up" {

	} else if action == "ui_down" {

	}
}

func containInSight(srcx Float, srcy Float, action string) {
	// if action == "ui_left" {
	// 	srcx - sight_value
	// } else if action == "ui_right" {
	// 	srcx + sight_value
	// } else if action == "ui_up" {
	// 	srcy - sight_value
	// } else if action == "ui_down" {
	// 	srcy + sight_value
	// }
}

func handleCommand(buf []byte, buf_len int, client_addr net.Addr) {
	buffstr := string(buf[:])
	headers := SpliteMsg(buffstr)
	userid := headers[0]
	command := headers[1]
	if command == "ping" {
		fmt.Println("ping pong")
		var msg []byte = []byte{'p', 'o', 'n', 'g'}
		msg_len := len(msg)
		unicast(userid, msg, msg_len)
		return
	} else if command == "join" {
		handleJoin(userid, client_addr)
	} else if command == "move" {
		action := headers[2]
		//delta_time := headers[3]
		//speed := headers[4]
		pos := headers[5]
		last_position[userid] = pos
		pos_v2, _ := GetV2Pos(pos)
		handleMove(pos_v2.X, pos_v2.Y, action)
		// fmt.Pintln(userid + " last pos : " + pos)
	}
	fmt.Println("new : ", buf_len)
	msg_queue_ch <- NewMsgBuff(buf, uint32(buf_len))
}

func GetV2Pos(str string) (Vector2, error) {
	strings.Trim(str, "()")
	tok := ", "
	p := strings.Index(str, tok)
	if p == -1 {
		return Vector2{}, fmt.Errorf("")
	}
	x, _ := strconv.ParseFloat(str[:p], 32)
	y, _ := strconv.ParseFloat(str[p+len(tok):], 32)
	v := Vector2{Float(x), Float(y)}
	return v, nil
}
