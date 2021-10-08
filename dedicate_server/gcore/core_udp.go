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

// var server net.PacketConn

// const MAX_BUFFSIZE = 1500

// func Send(buf []byte, buf_len int, to net.Addr) (int, error) {
// 	sent := 0
// 	leftsize := buf_len
// 	for leftsize >= MAX_BUFFSIZE {
// 		n, err := server.WriteTo(buf[sent:sent+MAX_BUFFSIZE], to)
// 		if err != nil {
// 			fmt.Println("Send error " + to.String() + " : " + err.Error())
// 			return n, err
// 		}
// 		sent += MAX_BUFFSIZE
// 		leftsize -= MAX_BUFFSIZE
// 	}
// 	if leftsize > 0 {
// 		n, err := server.WriteTo(buf[sent:sent+leftsize], to)
// 		if err != nil {
// 			fmt.Println("Send error " + to.String() + " : " + err.Error())
// 			return n, err
// 		}
// 		sent += n
// 	}
// 	return sent, nil
// }

// func broadcast(buf []byte, buf_len int) {
// 	for _, player := range GetWorld().Players {
// 		// Test
// 		_, err := server.WriteTo(buf[:buf_len], player.Addr) //
// 		// _, err := Send(buf[:buf_len], buf_len, player.Addr)
// 		if err != nil {
// 			fmt.Println("broadcast error " + player.NickName + ": " + err.Error())
// 			delete(GetWorld().Players, player.NickName)
// 		}
// 	}
// }

// func unicast(userid string, buf []byte, buf_len int) {
// 	player := GetWorld().Players[userid]
// 	_, err := Send(buf[:buf_len], buf_len, player.Addr)
// 	if err != nil {
// 		fmt.Println("unicast error " + player.NickName + ": " + err.Error())
// 		delete(GetWorld().Players, player.NickName)
// 	}
// }

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
	//go SyncAllObjects(200)
	go SyncAllObjects(200)

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
