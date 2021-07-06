package server

import (
	"fmt"
	"io"
	"log"
)

type Handler interface {
	BeforeHandle(Request)
	Handle(Request)
	AfterHandle(Request)
}

//type BaseHandler struct {}
//
//func (b *BaseHandler) BeforeHandle(r Request) {
//
//}
//
//func (b *BaseHandler) Handle(r Request) {
//
//}
//
//func (b *BaseHandler) AfterHandle(r Request) {
//
//}

type EchoHandler struct {}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (e *EchoHandler) BeforeHandle(r Request) {
	fmt.Println("before func")
}

func (e *EchoHandler) Handle(r Request) {
	for {
		data, err := r.Conn().Receive()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			return
		}

		_, err = r.Conn().Send(data, r.MsgType())
		if err != nil {
			log.Println(err)
			return
		}
	}

}

func (e *EchoHandler) AfterHandle(r Request) {
	fmt.Println("after func")
}
