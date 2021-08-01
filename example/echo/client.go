package main

import (
	"fmt"
	"log"
	"time"
	"zinx/client"
	"zinx/server"
)

func main() {
	conn := client.NewClientWithTCP("localhost", "8080")
	// 发送心跳包
	go conn.SendHeartbeat()

	msgs := []string{"666", "456", "789"}

	for _, msg := range msgs {
		msg := server.NewMessage([]byte(msg), server.EchoMsg)
		n, err := conn.Send(msg)
		log.Printf("send %d bytes", n)
		if err != nil {
			log.Fatalln(err)
		}

		receive, err := conn.Receive()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(string(receive.Data()))
		time.Sleep(time.Second * 3)
	}
}
