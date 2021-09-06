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

const LOCKSTEP_CNT = 15                                          // ?? miliseconds per onetime
var client_addrs map[string]net.Addr = make(map[string]net.Addr) // k : user_id, v : address
var entered_clients map[string]string = make(map[string]string)
var last_position map[string]string = make(map[string]string)
var server net.PacketConn

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

func GetNowTimeMili() int64 {
	return time.Now().UnixNano() / 1000000
}

func ExecLockstep() {
	sending_buffer := make([]byte, MAX_BUFFSIZE)
	var sending_size uint32 = 0
	var last_timestamp int64 = GetNowTimeMili()
	var now_timestamp int64
	for {
		now_timestamp = GetNowTimeMili()
		msg := <-msg_queue_ch
		copy(sending_buffer[sending_size:], msg.data)
		sending_size += msg.size
		for (now_timestamp - last_timestamp) <= LOCKSTEP_CNT {
			now_timestamp = GetNowTimeMili()
		}
		last_timestamp = now_timestamp
		//fmt.Println("duration: ", (now_mili - last_mili))
		//나머지
		for sending_size < DGRAM_SIZE {
			select {
			case msg := <-msg_queue_ch:
				copy(sending_buffer[sending_size:], msg.data)
				sending_size += msg.size
			default:
				broadcast(sending_buffer, sending_size)
				sending_size = 0
				return
			}
		}
		broadcast(sending_buffer, sending_size)
		sending_size = 0
	}
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
