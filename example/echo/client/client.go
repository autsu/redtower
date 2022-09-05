package main

import (
	"context"
	"github.com/autsu/redtower/client"
	"github.com/autsu/redtower/example"
	"github.com/autsu/redtower/server"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Ltime)
	ctx := context.Background()
	cli, err := client.NewClientWithTCP("localhost", "7788").Init(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	msgs := []string{"666", "456", "789", "999", "zxzxz", "12435", "dsfsd"}

	for i := 0; i < len(msgs); i++ {
		msg := server.NewMessage([]byte(msgs[i]), example.EchoMsg)
		_, err := cli.Send(msg)
		//log.Printf("send %d bytes", n)
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
