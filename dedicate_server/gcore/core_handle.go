package gcore

import (
	"container/list"
	"fmt"
	. "gnet"
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
			Broadcast(sending_buffer, sending_size)
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
		peers := GetPeers()
		for name, _ := range *peers {
			for i := 0; i < 100; i++ {
				_, err := Unicast(name, []byte(msg), len(msg))
				if err != nil {
					LeavePeer(name)
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
			Unicast(player.NickName, []byte(msg), len(msg))
		}
	}
}

func SyncAllSzObjects(tick_mili time.Duration) {
	ticker := time.NewTicker(time.Millisecond * tick_mili)
	for _ = range ticker.C {
		var founds_ch chan *list.List
		go GetWorld().object_tree.GetAllObjectsToCh(founds_ch)
		// if founds_ch == nil {
		// 	continue
		// }
		for founds := range founds_ch {
			var msgarr []byte = make([]byte, 1024)
			pack := NewSyncPacket(TYPE_HEADER_SYNC, msgarr, int32(len(msgarr)))
			//msgarr = append(msgarr, "noti;objects;"...)
			for e := founds.Front(); e != nil; e = e.Next() {
				obj := e.Value.(*GObject)
				if obj != nil {
					data, _ := obj.Serialize()
					msgarr = append(msgarr, data...)
					msgarr = append(msgarr, "@"...)
				}
			}
			msgarr = append(msgarr, ";m;"...)
			pack.Data = msgarr
			packet_que <- pack
		}

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
			_, err := Unicast(player.NickName, syncmsg_buff, syncmsg_len)
			// _, err := server.WriteTo(syncmsg_buff[:syncmsg_len], client_addr)
			if err != nil {
				log.Fatal(err)
				return false
			}
		}
		fmt.Println("enter client : " + client_addr.String() + ", " + nickname)
		GetWorld().AddPlayer(nickname, client_addr, pos)
		GetWorld().Players[nickname] = NewPlayer(0, nickname, client_addr, pos, DEFAULT_COLISION_RADIUS)
	}
	return true
}

func handleMove(player *Player, action string) {
	return
}

func handleCommand(buf []byte, buf_len int, client_addr net.Addr) {
	packet := ParsePacketHeader(buf)
	if packet.HeaderType == TYPE_HEADER_CMD {
		if packet.CommandType == TYPE_COMMAND_ENTER {
			userid, pos_v2 := ParseCommandData(packet.Data)
			enterClient(userid, client_addr, pos_v2)
			fmt.Println("enter ", userid)
			// screen_size := header[2]
		} else if packet.CommandType == TYPE_COMMAND_MOVE {
			// action := header[2]
			// handleMove(player, action)
			return
		}
		msg_queue_ch <- NewMsgBuff(buf, buf_len)

	} else if packet.HeaderType == TYPE_HEADER_SYNC {
		userid, pos_v2 := ParseSyncData(packet.Data)
		objname := header[2]
		pos := header[3]
		pos_v2, _ := posStr2V2(pos)
		GetWorld().AddObject(&GObject{Name: objname, Pos: pos_v2})
	}
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
