package gnet

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	. "dedicate_server/gcore"
)

var msg_queue_ch = make(chan *MsgBuff)

const FPS = 60
const LOCKSTEP_CNT = 200 // ?? miliseconds per onetime

func GetNowTimeMili() int64 {
	return time.Now().UnixNano() / 1000000
}
func ExecLockstep() {
	sending_buffer := make([]byte, MAX_BUFFSIZE)
	var sending_size uint32 = 0
	var last_timestamp int64 = GetNowTimeMili()
	var now_timestamp int64
	var duration int64
	for {
		msg := <-msg_queue_ch
		copy(sending_buffer[sending_size:], msg.data)
		sending_size += msg.size
		now_timestamp = GetNowTimeMili()
		duration = (now_timestamp - last_timestamp)
		if duration <= LOCKSTEP_CNT {
			continue
		}
		last_timestamp = now_timestamp
		duration = 0
		//나머지
		for (sending_size != 0) && (sending_size < DGRAM_SIZE) {
			select {
			case msg := <-msg_queue_ch:
				copy(sending_buffer[sending_size:], msg.data)
				sending_size += msg.size
			default:
				broadcast(sending_buffer, sending_size)
				sending_size = 0
			}
		}
		if sending_size > 0 {
			broadcast(sending_buffer, sending_size)
			sending_size = 0
		}
	}
}

func enterClient(userid string, client_addr net.Addr, pos Vector2) bool {
	_, exists := GetWorld().Players[userid]
	if exists == false {
		for _, player := range GetWorld().Players {
			if userid == player.Uid {
				continue
			}
			syncmsg := player.Uid + ";sync;" + GetWorld().Players[player.Uid].GetPositionStr() + ";m;"
			syncmsg_len := len(syncmsg)
			syncmsg_buff := make([]byte, syncmsg_len)
			copy(syncmsg_buff, syncmsg)
			_, err := server.WriteTo(syncmsg_buff[:syncmsg_len], client_addr)
			if err != nil {
				log.Fatal(err)
				return false
			}
		}
		fmt.Println("enter client : " + client_addr.String() + ", " + userid)
		GetWorld().Players[userid] = NewPlayer(userid, client_addr, pos, DEFAULT_COLISION_RADIUS)
	}
	return true
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
	} else if command == "enter" {
		player_pos := headers[2]
		pos_v2, _ := posStr2V2(player_pos)
		enterClient(userid, client_addr, pos_v2)
		// screen_size := headers[2]
	} else if command == "move" {
		action := headers[2]
		//delta_time := headers[3]
		//speed := headers[4]
		pos := headers[5]
		pos_v2, _ := posStr2V2(pos)
		player := GetWorld().Players[userid]
		if player != nil {
			player.UpdatePos(pos_v2)
		} else {
			fmt.Println("nil player " + userid)
		}
		handleMove(pos_v2.X, pos_v2.Y, action)
		// fmt.Pintln(userid + " last pos : " + pos)
	}
	msg_queue_ch <- NewMsgBuff(buf, uint32(buf_len))
}

// "(40, 40)" -> x:40, y:40 Float Vector2
func posStr2V2(str string) (Vector2, error) {
	str = strings.Trim(str, "()")
	tok := ", "
	p := strings.Index(str, tok)
	if p == -1 {
		return Vector2{}, fmt.Errorf("invalid value " + str)
	}
	x, _ := strconv.ParseFloat(str[:p], 32)
	y, _ := strconv.ParseFloat(str[p+len(tok):], 32)
	v := Vector2{Float(x), Float(y)}
	return v, nil
}
