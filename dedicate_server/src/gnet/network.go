package gnet

import (
	"dedicate_server/gcore"
	"fmt"
	"log"
	"net"
	"strconv"
)

// * Pakcet format
// * requirements : user_id;command;
// * user_id;command;action;delta-time;

var server net.PacketConn

const DGRAM_SIZE = 1400
const MAX_BUFFSIZE = 1500

func broadcast(buf []byte, buf_len int32) {
	for _, player := range gcore.GetWorld().Players {
		_, err := server.WriteTo(buf[:buf_len], player.Addr)
		if err != nil {
			fmt.Println("broadcast error " + player.NickName + ": " + err.Error())
			delete(gcore.GetWorld().Players, player.NickName)
		}
	}
}

func unicast(userid string, buf []byte, buf_len int) {
	player := gcore.GetWorld().Players[userid]
	_, err := server.WriteTo(buf[:buf_len], player.Addr)
	if err != nil {
		fmt.Println("unicast error " + player.NickName + ": " + err.Error())
		delete(gcore.GetWorld().Players, player.NickName)
	}
}

func RunUDPServer(server_ip string, server_port int) net.Addr {
	serv_addr := server_ip + ":" + strconv.Itoa(server_port)
	var err error
	server, err = net.ListenPacket("udp", serv_addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server address: ", server.LocalAddr().String())

	go ExecLockstep()

	// 200ms 마다 오브젝트 동기화
	go SyncObjects(200)

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
