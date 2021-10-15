package gcore

import (
	"encoding/json"
	"fmt"
	"gnet"
	"io"
	"log"
	"net"
	"strconv"
)

// * Pakcet format
// * requirements : user_id;command;
// * user_id;command;action;delta-time;

func RunUdpServer(server_ip string, server_port int) net.Addr {
	serv_addr := server_ip + ":" + strconv.Itoa(server_port)
	// var err error
	server, err := net.ListenPacket("udp", serv_addr)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
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

func RunTcpServer(server_ip string, server_port int) bool {
	// TODO : you need change it to async read
	go func() {
		serv_addr := server_ip + ":" + strconv.Itoa(server_port)
		listen_conn, err := net.Listen("tcp", serv_addr)
		if err != nil {
			log.Fatal(err)
		}
		defer listen_conn.Close()
		fmt.Println("server address: ", listen_conn.Addr().String())
		for {
			conn, err := listen_conn.Accept()
			if nil != err {
				log.Println(err)
				continue
			}
			defer conn.Close()

			recvbuff := make([]byte, 4096) // receive buffer: 4kB
			if n, _ := conn.Read(recvbuff); n > 0 {
				packet := gnet.ParsePacketHeader(recvbuff)
				if packet.HeaderType == gnet.TYPE_HEADER_CMD {
					if packet.Command == gnet.TYPE_COMMAND_ENTER {
						doc := packet.GetData()
						var data map[string]interface{}
						data_len := n - gnet.HEADER_LENGTH
						err := json.Unmarshal([]byte(doc[:data_len]), &data)
						if err == nil {
							fmt.Println("joined, "+data["usr_id"], data["usr_pos"])
							usr_id := data["usr_id"].(string)
							pos := data["usr_pos"].(string)
							v2, _ := gnet.PosStr2V2(pos)
							handleEnterClient(usr_id, conn.RemoteAddr(), v2)
						}
					}
				}

			}

			go func() {
				for {
					//auth
					n, err := conn.Read(recvbuff)
					if nil != err {
						if io.EOF == err {
							log.Printf("connection is closed from client; %v", conn.RemoteAddr().String())
							return
						}
						log.Printf("fail to receive data; err: %v", err)
						return
					}
					if 0 < n {
						data := recvbuff[:n]
						log.Println(string(data))
					}
					handleCommand(recvbuff, n, conn.LocalAddr())
				}
			}()
		}
	}()

	go ExecLockstep()
	// 200ms 마다 오브젝트 동기화
	go SyncAllSzObjects(200)

	return true
}
