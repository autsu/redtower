package main

import (
	"bufio"
	"bytes"
	"context"
	"github.com/autsu/redtower/example"
	"github.com/autsu/redtower/server"
	"log"
	"net/http"
)

func main() {
	s := server.NewTCPServer("localhost", "7788", "test1",
		func(req *server.Request) error {
			typ := req.MsgType()
			switch typ {
			case example.HTTPMsg:
				handleHTTP(req)
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

func handleHTTP(req *server.Request) {
	bsr := bytes.NewReader(req.Data())
	br := bufio.NewReader(bsr)

	r, err := http.ReadRequest(br)
	if err != nil {
		log.Println(err)
		msg := server.NewErrorMessage([]byte(err.Error()))
		req.Conn().Send(msg)
		return
	}
	log.Printf("%+v", r)

	r.ParseForm()
	form := r.Form.Encode()
	msg := server.NewMessage([]byte(form), example.HTTPMsg)
	req.Conn().Send(msg)
}
