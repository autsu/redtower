package server

import (
	"fmt"
	"log"
)

type Handler interface {
	BeforeHandle(*Request)
	Handle(*Request)
	AfterHandle(*Request)
}

type EchoHandler struct {}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (e *EchoHandler) BeforeHandle(r *Request) {
	fmt.Println("before func")
}

func (e *EchoHandler) Handle(r *Request) {
	// 不需要再读了，conn.ReceiveAndHandler() 已经将数据读取出来了，
	// 读取的数据已经保存在 r.Data() 中，这里如果再次读取会无限阻塞

	//data, err := r.Conn().Receive()
	//if err != nil {
	//	if err == io.EOF {
	//		log.Println("EOF")
	//		break
	//	}
	//	log.Println(err)
	//	return
	//}
	//log.Println("receive data: ", data)

	data := r.Data()

	msg := NewMessage(data, r.MsgType())
	_, err := r.Conn().Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
}

func (e *EchoHandler) AfterHandle(r *Request) {
	fmt.Println("after func")
}
