package gcore

import (
	"fmt"
	. "libgnet/gnet"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

const circle_radius int = 200

var msg_queue_ch = make(chan *MsgBuff)
var packet_que = make(chan *SyncPacket)

const FPS = 60
const LOCKSTEP_CNT = 200 // ?? miliseconds per onetime
const DGRAM_SIZE = 1400

func GetNowTimeMili() int64 {
	return time.Now().UnixNano() / 1000000
}

func ExecLockstep() {
	sending_buffer := make([]byte, 65535)
	var sending_size int = 0
	var last_timestamp int64 = GetNowTimeMili()
	var now_timestamp int64
	var duration int64
	for {
		msg := <-msg_queue_ch
		copy(sending_buffer[sending_size:], msg.Data)
		sending_size += msg.Size
		//나머지
		for (sending_size != 0) && (sending_size < DGRAM_SIZE) {
			now_timestamp = GetNowTimeMili()
			duration = (now_timestamp - last_timestamp)
			if duration >= LOCKSTEP_CNT {
				break
			}
			select {
			case msg := <-msg_queue_ch:
				copy(sending_buffer[sending_size:], msg.Data)
				sending_size += msg.Size
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

func TEST_SyncObjects(tick_mili time.Duration) {
	min := -10
	max := 10
	rand.Seed(time.Now().UnixNano())
	ticker := time.NewTicker(time.Millisecond * tick_mili)
	for _ = range ticker.C {
		founds := GetWorld().GetAllObjects()
		if founds == nil {
			continue
		}
		var msg string = "noti;objects;"
		for e := founds.Front(); e != nil; e = e.Next() {
			detected_obj := e.Value.(*GObject)
			if detected_obj != nil {
				testtrigger := rand.Intn(max-min) + min
				x, y := (detected_obj.Pos.X + testtrigger), (detected_obj.Pos.Y + testtrigger + 1)
				msg += detected_obj.Name + "_" + ToPosString(x, y) + "@"
			}
		}
		msg += ";m;"
		for _, player := range GetWorld().Players {
			for i := 0; i < 100; i++ {
				_, err := server.WriteTo([]byte(msg)[:len(msg)], player.Addr)
				if err != nil {
					fmt.Println("broadcast error " + player.NickName + ": " + err.Error())
					delete(GetWorld().Players, player.NickName)
				}
			}
			break
		}
	}
}

func SyncNearObjects(tick_mili time.Duration) {
	ticker := time.NewTicker(time.Millisecond * tick_mili)
	for _ = range ticker.C {
		for _, player := range GetWorld().Players {
			founds := GetWorld().Nearest(player)
			if founds == nil {
				fmt.Println("not found")
				continue
			}
			var msg string = "noti;objects;"
			for e := founds.Front(); e != nil; e = e.Next() {
				detected_obj := e.Value.(*GObject)
				if detected_obj != nil {
					x, y := detected_obj.Pos.X, detected_obj.Pos.Y
					msg += detected_obj.Name + "_" + ToPosString(x, y) + "@"
				}
			}
			msg += ";m;"
			Send([]byte(msg), len(msg), player.Addr)
		}
	}
}

func SyncAllSzObjects(tick_mili time.Duration) {
	ticker := time.NewTicker(time.Millisecond * tick_mili)
	for _ = range ticker.C {
		founds := GetWorld().GetAllObjects()
		if founds == nil {
			continue
		}
		var msgarr []byte = make([]byte, 1024)
		msgarr = append(msgarr, "noti;objects;"...)
		for e := founds.Front(); e != nil; e = e.Next() {
			obj := e.Value.(*GObject)
			if obj != nil {
				data, _ := obj.Serialize()
				msgarr = append(msgarr, data...)
				msgarr = append(msgarr, "@"...)
			}
		}
		msgarr = append(msgarr, ";m;"...)
		packet_que <- NewSyncPacket(HEADERTYPE_SYNC, msgarr, int32(len(msgarr)))
	}
}

func SyncAllObjects(tick_mili time.Duration) {
	ticker := time.NewTicker(time.Millisecond * tick_mili)
	for _ = range ticker.C {
		founds := GetWorld().GetAllObjects()
		if founds == nil {
			continue
		}
		var msg string = "noti;objects;"
		for e := founds.Front(); e != nil; e = e.Next() {
			obj := e.Value.(*GObject)
			if obj != nil {
				msg += obj.Name + "_" + ToPosString(obj.Pos.X, obj.Pos.Y) + "@"
			}
		}
		msg += ";m;"
		msg_queue_ch <- NewMsgBuff([]byte(msg), len(msg))
	}

}
func enterClient(nickname string, client_addr net.Addr, pos Vector2) bool {
	_, exists := GetWorld().Players[nickname]
	if exists == false {
		for _, player := range GetWorld().Players {
			if nickname == player.NickName {
				continue
			}
			syncmsg := player.NickName + ";sync;" + GetWorld().Players[player.NickName].GetPositionStr() + ";m;"
			syncmsg_len := len(syncmsg)
			syncmsg_buff := make([]byte, syncmsg_len)
			copy(syncmsg_buff, syncmsg)
			_, err := server.WriteTo(syncmsg_buff[:syncmsg_len], client_addr)
			if err != nil {
				log.Fatal(err)
				return false
			}
		}
		fmt.Println("enter client : " + client_addr.String() + ", " + nickname)
		GetWorld().Players[nickname] = NewPlayer(0, nickname, client_addr, pos, DEFAULT_COLISION_RADIUS)
	}
	return true
}

func handleMove(player *Player, action string) {
	// if player != nil {
	// 	player.UpdatePos(pos_v2)
	// } else {
	// 	fmt.Println("nil player " + userid)
	// }
	return
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

		//TODO: Player객체를 업데이트 하고 오브젝트 위치 변경에 따라 objects_tree에도 업데이트가 되어야함
		player := GetWorld().Players[userid]
		if player != nil {
			player.UpdatePos(pos_v2)
		} else {
			fmt.Println("nil player " + userid)
		}
		handleMove(player, action)
		return
	} else if command == "noti" {
		if userid == "obj" {
			objname := header[2]
			pos := header[3]
			pos_v2, _ := posStr2V2(pos)
			GetWorld().AddObject(&GObject{Name: objname, Pos: pos_v2})
		}
	}
	msg_queue_ch <- NewMsgBuff(buf, buf_len)
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
func ToPosString(x int, y int) string {
	return "(" + strconv.Itoa(int(x)) + ", " + strconv.Itoa(int(y)) + ")"
}
func v2Str(v Vector2) string {
	return "(" + strconv.Itoa(int(v.X)) + ", " + strconv.Itoa(int(v.Y)) + ")"
}
