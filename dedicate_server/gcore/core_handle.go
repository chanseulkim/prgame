package gcore

import (
	"container/list"
	"fmt"
	. "gnet"
	"math/rand"
	"net"
	"time"
)

const circle_radius int = 200

var msg_queue_ch = make(chan *MsgBuff)
var packet_que = make(chan *GPacket)

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
		packet := <-packet_que
		copy(sending_buffer[sending_size:], packet.GetData())
		sending_size += int(packet.GetDataLength())
		//나머지
		for (sending_size != 0) && (sending_size < DGRAM_SIZE) {
			now_timestamp = GetNowTimeMili()
			duration = (now_timestamp - last_timestamp)
			if duration >= LOCKSTEP_CNT {
				break
			}
			select {
			case packet := <-packet_que:
				copy(sending_buffer[sending_size:], packet.GetData())
				sending_size += int(packet.GetDataLength())
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
		var founds_ch chan *list.List = make(chan *list.List)
		go GetWorld().object_tree.GetAllObjectsToCh(founds_ch)
		// if founds_ch == nil {
		// 	continue
		// }
		var allobjs *list.List = list.New()
		for founds := range founds_ch {
			allobjs.PushBackList(founds)
		}
		pack := NewSyncPacket(TYPE_PACKET_WHOLE, allobjs)
		packet_que <- pack
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

func handleCommand(buf []byte, buf_len int, client_addr net.Addr) {
	packet := ParsePacketHeader(buf)
	if packet.HeaderType == TYPE_HEADER_CMD {
		if packet.Command == TYPE_COMMAND_ENTER {
			userid, pos_v2 := ParseCommandData(packet.GetData())
			handleEnterClient(userid, client_addr, pos_v2)
			fmt.Println("enter ", userid)
			// screen_size := header[2]
		} else if packet.Command == TYPE_COMMAND_MOVE {
			// action := header[2]
			// handleMove(player, action)
			return
		}
		msg_queue_ch <- NewMsgBuff(buf, buf_len)
	}
}
