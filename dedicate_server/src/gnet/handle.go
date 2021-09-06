package gnet

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	. "dedicate_server/gcore"
)

var colision_radius Float = 20.0
var sight_value Float = 20.0
var msg_queue_ch = make(chan *MsgBuff)

var avg_delay int64 = 0
var fps_val int64 = 0
var frame_count int = 0

const fps = 60

func checkFps(duration int64) {
	frame_count++
	if frame_count >= fps {
		avg_delay = fps_val / fps
		fmt.Println("avg_delay : ", avg_delay)
		frame_count = 0
		fps_val = 0
		return
	}
	fps_val += duration
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
		last_timestamp = GetNowTimeMili()
		for duration <= LOCKSTEP_CNT {
			now_timestamp = GetNowTimeMili()
			duration = (now_timestamp - last_timestamp)
		}
		checkFps(duration)
		duration = 0
		last_timestamp = now_timestamp
		//나머지
		for (sending_size != 0) && (sending_size < DGRAM_SIZE) {
			select {
			case msg := <-msg_queue_ch:
				copy(sending_buffer[sending_size:], msg.data)
				sending_size += msg.size
			default:
				broadcast(sending_buffer, sending_size)
				sending_size = 0
				break
			}
		}
		if sending_size > 0 {
			broadcast(sending_buffer, sending_size)
			sending_size = 0
		}
	}
}

func EnterClient(userid string, client_addr net.Addr, pos Vector2) {
	_, exists := GetWorld().Players[userid]
	if exists == false {
		for _, player := range GetWorld().Players {
			syncmsg := player.Uid + ";sync;" + GetWorld().Players[player.Uid].GetPositionStr() + ";m;"
			syncmsg_len := len(syncmsg)
			syncmsg_buff := make([]byte, syncmsg_len)
			copy(syncmsg_buff, syncmsg)
			_, err := server.WriteTo(syncmsg_buff[:syncmsg_len], client_addr)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("enter client : " + client_addr.String() + ", " + userid)
		//client_addrs[userid] = client_addr
		GetWorld().Players[userid] = NewPlayer(userid, client_addr, pos, colision_radius)
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
	} else if command == "enter" {
		pos := headers[5]
		pos_v2, _ := GetPosV2(pos)
		EnterClient(userid, client_addr, pos_v2)
	} else if command == "move" {
		action := headers[2]
		//delta_time := headers[3]
		//speed := headers[4]
		pos := headers[5]
		pos_v2, _ := GetPosV2(pos)
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

func GetPosV2(str string) (Vector2, error) {
	strings.Trim(str, "()")
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
