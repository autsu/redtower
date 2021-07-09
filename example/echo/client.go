package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"zinx/conf"
	"zinx/server"
)

func main() {
	tcpaddr, _ := net.ResolveTCPAddr("tcp", ":8080")
	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		log.Fatalln(err)
	}

	tcpConn := server.NewTCPConn(conn, nil, 1)
	msgs := []string{"666", "456", "789"}

	// 发送心跳包
	go func() {
		ticker := time.NewTicker(conf.SendHeartbeatTime)
		for {
			select {
			case <-ticker.C:
				heartbeat := server.NewMessage([]byte(""), server.HeartBeat)
				tcpConn.Send(heartbeat)
				//log.Println("send heartbreat to server")
			}
		}
	}()

	for _, msg := range msgs {
		msg := server.NewMessage([]byte(msg), server.EchoMsg)
		n, err := tcpConn.Send(msg)
		log.Printf("send %d bytes", n)
		if err != nil {
			log.Fatalln(err)
		}

		receive, err := tcpConn.Receive()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(string(receive.Data()))
		time.Sleep(time.Second * 3)
	}
}
