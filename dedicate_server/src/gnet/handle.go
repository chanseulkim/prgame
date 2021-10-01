package gnet

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	. "dedicate_server/gcore"
)

const circle_radius int = 200

var msg_queue_ch = make(chan *MsgBuff)

const FPS = 60
const LOCKSTEP_CNT = 200 // ?? miliseconds per onetime

func GetNowTimeMili() int64 {
	return time.Now().UnixNano() / 1000000
}

func ExecLockstep() {
	sending_buffer := make([]byte, MAX_BUFFSIZE)
	var sending_size int32 = 0
	var last_timestamp int64 = GetNowTimeMili()
	var now_timestamp int64
	var duration int64
	for {
		msg := <-msg_queue_ch
		copy(sending_buffer[sending_size:], msg.data)
		sending_size += msg.size
		//나머지
		for (sending_size != 0) && (sending_size < DGRAM_SIZE) {
			now_timestamp = GetNowTimeMili()
			duration = (now_timestamp - last_timestamp)
			if duration >= LOCKSTEP_CNT {
				break
			}
			select {
			case msg := <-msg_queue_ch:
				copy(sending_buffer[sending_size:], msg.data)
				sending_size += msg.size
			default:
				continue
			}
		}
		last_timestamp = now_timestamp
		if sending_size > 0 {
			broadcast(sending_buffer, sending_size)
			sending_size = 0
		}

	}
}
func SyncObjects(tick_mili time.Duration) {
	min := -10
	max := 10
	rand.Seed(time.Now().UnixNano())

	// var objque = make(chan *string)
	ticker := time.NewTicker(time.Millisecond * tick_mili)
	for _ = range ticker.C {
		founds := GetWorld().GetAllObjects()
		if founds == nil {
			continue
		}
		var msg string = "noti;objects;"
		for _, detected_obj := range founds {
			if detected_obj != nil {
				testtrigger := rand.Intn(max-min) + min
				x, y := (detected_obj.Pos.X + testtrigger), (detected_obj.Pos.Y + testtrigger + 1)
				msg += detected_obj.Name + "_" + v2Str(Vector2{X: int(x), Y: int(y)}) + "@"
			}
		}
		msg += ";m;"
		// objque <- &msg
		broadcast([]byte(msg), int32(len(msg)))
	}

}
func enterClient(userid string, client_addr net.Addr, pos Vector2) bool {
	_, exists := GetWorld().Players[userid]
	if exists == false {
		for _, player := range GetWorld().Players {
			if userid == player.UsrId {
				continue
			}
			syncmsg := player.UsrId + ";sync;" + GetWorld().Players[player.UsrId].GetPositionStr() + ";m;"
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

func handleMove(player *Player, action string) {
	return
	// TODO: 시야거리 적용
	// founds := GetWorld().Nearest(player)
	// var msg string
	// if founds == nil {
	// 	msg += "noti;objects;m;"
	// 	unicast(player.UsrId, []byte(msg), len(msg))
	// 	return
	// }
	// msg += "noti;objects;"
	// for _, detected_obj := range founds {
	// 	if detected_obj != nil {
	// 		//log.Printf("Found point: %s\n", detected_obj.Name)
	// 		x, y := detected_obj.Pos.X, detected_obj.Pos.Y
	// 		msg += detected_obj.Name + "_" + v2Str(Vector2{X: int(x), Y: int(y)}) + "@"
	// 	}
	// }
	// msg += ";m;"
	// msglen := len(msg)
	// unicast(player.UsrId, []byte(msg), msglen)
}

func handleCommand(buf []byte, buf_len int, client_addr net.Addr) {
	buffstr := string(buf[:])
	header := SpliteMsg(buffstr)
	userid := header[0]
	command := header[1]
	if command == "ping" {
		fmt.Println("ping pong")
		var msg []byte = []byte{'p', 'o', 'n', 'g'}
		msg_len := len(msg)
		unicast(userid, msg, msg_len)
		return
	} else if command == "enter" {
		player_pos := header[2]
		pos_v2, _ := posStr2V2(player_pos)
		enterClient(userid, client_addr, pos_v2)
		fmt.Println("enter ", userid)
		// screen_size := header[2]
	} else if command == "move" {
		action := header[2]
		//delta_time := header[3]
		//speed := header[4]
		pos := header[5]
		pos_v2, _ := posStr2V2(pos)
		player := GetWorld().Players[userid]
		if player != nil {
			player.UpdatePos(pos_v2)
		} else {
			fmt.Println("nil player " + userid)
		}
		handleMove(player, action)

	} else if command == "noti" {
		if userid == "obj" {
			objname := header[2]
			pos := header[3]
			pos_v2, _ := posStr2V2(pos)
			GetWorld().AddObject(&GObject{Name: objname, Pos: pos_v2})
		}
	}
	msg_queue_ch <- NewMsgBuff(buf, int32(buf_len))
}

// "(40, 40)" -> x:40, y:40 int Vector2
func posStr2V2(str string) (Vector2, error) {
	str = strings.Trim(str, "()")
	tok := ", "
	p := strings.Index(str, tok)
	if p == -1 {
		return Vector2{}, fmt.Errorf("invalid value " + str)
	}
	x, _ := strconv.ParseFloat(str[:p], 32)
	y, _ := strconv.ParseFloat(str[p+len(tok):], 32)
	v := Vector2{int(x), int(y)}
	return v, nil
}
func v2Str(v Vector2) string {
	return "(" + strconv.Itoa(int(v.X)) + ", " + strconv.Itoa(int(v.Y)) + ")"
}
