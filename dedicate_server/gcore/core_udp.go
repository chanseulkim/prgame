package gcore

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

// * Pakcet format
// * requirements : user_id;command;
// * user_id;command;action;delta-time;

func RunUDPServer(server_ip string, server_port int) net.Addr {
	serv_addr := server_ip + ":" + strconv.Itoa(server_port)
	// var err error
	server, err := net.ListenPacket("udp", serv_addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("server address: ", server.LocalAddr().String())

	go ExecLockstep()

	// 200ms 마다 오브젝트 동기화
	go SyncAllSzObjects(200)

	recvbuff := make([]byte, 65535)
	go func() {
		for {
			n, clientAddress, err := server.ReadFrom(recvbuff)
			if (n == 0) || (err != nil) {
				fmt.Println("buffer size is 0...")
			}
			handleCommand(recvbuff, n, clientAddress)
		}
	}()
	return server.LocalAddr()
}
