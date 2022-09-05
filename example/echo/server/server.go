package main

import (
	"context"
	"github.com/autsu/redtower/example"
	"github.com/autsu/redtower/server"
	"log"
)

func main() {
	s := server.NewTCPServer("localhost", "7788", "test1",
		func(req *server.Request) error {
			typ := req.MsgType()
			switch typ {
			case example.EchoMsg:
				handleEcho(req)
			case server.ErrorMsg:
				log.Printf("error: %v\n", string(req.Data()))
			default:
				log.Println("Unknown msg typ")
				log.Println(req.Data())
			}
			return nil
		})
	if err := s.Run(context.Background()); err != nil {
		panic(err)
	}
}

func handleEcho(req *server.Request) error {
	data := req.Data()
	msg := server.NewMessage(data, example.EchoMsg)
	if _, err := req.Conn().Send(msg); err != nil {
		return err
	}
	return nil
}
