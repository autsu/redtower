package main

import (
	"context"
	"log"
	"time"
	"zinx/client"
	"zinx/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Ltime)
	cli := client.NewClientWithTCP("localhost", "8080")
	cli.StartHeartbeat(context.Background())

	msgs := []string{"666", "456", "789", "999", "zxzxz", "12435", "dsfsd"}

	for i := 0; i < len(msgs); i++ {
		msg := server.NewMessage([]byte(msgs[i]), server.EchoMsg)
		n, err := cli.Send(msg)
		log.Printf("send %d bytes", n)
		if err != nil {
			log.Println("send error: ", err)
			continue
		}

		receive, err := cli.Receive()
		if err != nil {
			log.Println("recv error: ", err)
			continue
		}

		log.Println(string(receive.Data()))
		time.Sleep(time.Second * 3)

		if i == len(msgs)-1 {
			i = 0
		}
	}
}
