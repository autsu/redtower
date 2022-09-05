package main

import (
	"context"
	"fmt"
	"github.com/autsu/redtower/client"
	"github.com/autsu/redtower/example"
	"github.com/autsu/redtower/server"
	"log"
)

func main() {
	ctx := context.Background()
	conn, err := client.
		NewClientWithTCP("localhost", "7788").
		Init(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	data := `POST /?123=456 HTTP/1.1

User-Agent: PostmanRuntime/7.28.1
Accept: */*
Postman-Token: e1e457b7-d713-443d-8022-04a2d1d9697a
Host: 127.0.0.1:8080
Accept-Encoding: gzip, deflate, br
Connection: keep-alive
Content-Length: 0`

	msg := server.NewMessage([]byte(data), example.HTTPMsg)
	conn.Send(msg)

	msg, err = conn.Receive()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(string(msg.Data()))
}
