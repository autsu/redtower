package main

import (
	"context"
	"github.com/autsu/redtower/server"
	"google.golang.org/protobuf/proto"
	"log"
	"testdata/net"
	"testdata/pbfile"
)

func main() {
	t := server.NewTCPServer("localhost", "7788", "testserver",
		func(req *server.Request) error {
			typ := req.MsgType()
			switch typ {
			case net.ProtoMsg:
				if err := handleProto(req); err != nil {
					return err
				}
			case server.ErrorMsg:
				log.Printf("error: %v\n", string(req.Data()))
			default:
				log.Println("Unknown msg typ")
				log.Println(req.Data())
			}
			return nil
		})
	if err := t.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func handleProto(req *server.Request) error {
	var m pbfile.Message
	if err := proto.Unmarshal(req.Data(), &m); err != nil {
		return err
	}
	log.Printf("%+v\n", m.Uid)
	return nil
}
