package gnet

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

// * Pakcet format
// * requirements : user_id;command;
// * user_id;command;action;delta-time;
// var client_addrs *list.List = list.New()

const LOCKSTEP_CNT = 15                                          // miliseconds
var client_addrs map[string]net.Addr = make(map[string]net.Addr) // k : user_id, v : address
var joined_clients map[string]string = make(map[string]string)
var last_position map[string]string = make(map[string]string)
var server net.PacketConn

var msg_queue_ch = make(chan *MsgBuff)

const DGRAM_SIZE = 1400
const MAX_BUFFSIZE = 1500

func broadcast(buf []byte, buf_len uint32) {
	for uid, addr := range client_addrs {
		_, err := server.WriteTo(buf[:buf_len], addr)
		if err != nil {
			delete(client_addrs, uid)
		}
	}
}

func unicast(userid string, buf []byte, buf_len int) {
	addr := client_addrs[userid]
	_, err := server.WriteTo(buf[:buf_len], addr)
	if err != nil {
		delete(client_addrs, userid)
	}
}

type MsgBuff struct {
	data     []byte
	size     uint32
	capacity uint32
}

func NewMsgBuff(data []byte, data_size uint32) *MsgBuff {
	var buff MsgBuff
	buff.data = make([]byte, data_size)
	copy(buff.data, data)
	buff.size = data_size
	return &buff
}

func (self *MsgBuff) Write(data []byte, data_size uint32) {
	copy(self.data[self.size:], data)
	self.size += data_size
}

func GetNowTimeMili() int64 {
	return time.Now().UnixNano() / 1000000
}

var forcheck int = 0

func ExecLockstep() {
	sending_buffer := make([]byte, MAX_BUFFSIZE)
	var sending_size uint32 = 0
	var last_mili int64 = GetNowTimeMili()
	var now_mili int64
	for {
		now_mili = GetNowTimeMili()
		msg := <-msg_queue_ch
		copy(sending_buffer[sending_size:], msg.data)
		sending_size += msg.size
		for (now_mili - last_mili) <= LOCKSTEP_CNT {
			now_mili = GetNowTimeMili()
		}
		fmt.Println("duration: ", (now_mili - last_mili))
		last_mili = now_mili
		//나머지
		empty_ch := false
		for (sending_size < DGRAM_SIZE) && (empty_ch == false) {
			select {
			case msg := <-msg_queue_ch:
				copy(sending_buffer[sending_size:], msg.data)
				sending_size += msg.size
			default:
				empty_ch = true
			}
		}
		broadcast(sending_buffer, sending_size)
		sending_size = 0
	}
}

func handleJoin(userid string, client_addr net.Addr) {
	_, exists := joined_clients[client_addr.String()]
	if exists == false {
		for _, usrid := range joined_clients {
			syncmsg := usrid + ";sync;" + last_position[usrid] + ";"
			syncmsg_len := len(syncmsg)
			syncmsg_buff := make([]byte, syncmsg_len)
			copy(syncmsg_buff, syncmsg)
			_, err := server.WriteTo(syncmsg_buff[:syncmsg_len], client_addr)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println("joined client : " + client_addr.String() + ", " + userid)
		client_addrs[userid] = client_addr
		joined_clients[client_addr.String()] = userid
	}
}
func handleCommand(buf []byte, buf_len int, client_addr net.Addr) {
	buffstr := string(buf[:])
	headers := ParseMsg(buffstr)
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
		// action := headers[2]
		// if action == "ui_left" {
		// 	//
		// }

		//delta_time := headers[3]
		//speed := headers[4]
		pos := headers[5]
		last_position[userid] = pos
		// fmt.Println(userid + " last pos : " + pos)
	}
	fmt.Println("new : ", buf_len)
	msg_queue_ch <- NewMsgBuff(buf, uint32(buf_len))
}

func MakeUDPServer(server_ip string, server_port int) net.Addr {
	serv_addr := server_ip + ":" + strconv.Itoa(server_port)
	var err error
	server, err = net.ListenPacket("udp", serv_addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server address: ", server.LocalAddr().String())
	go ExecLockstep()
	go func() {
		for {
			buf := make([]byte, MAX_BUFFSIZE)
			n, clientAddress, err := server.ReadFrom(buf)
			if (n == 0) || (err != nil) {
				fmt.Println("buffer size is 0...")
			}
			handleCommand(buf, n, clientAddress)
		}
	}()
	return server.LocalAddr()
}
