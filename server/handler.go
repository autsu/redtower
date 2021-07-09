package server

import (
	"log"
)

type Handler interface {
	BeforeHandle(*Request)
	// Handle 不需要再从 conn 中读取了，在该方法的上层调用 conn.Handler() 中
	// 已经将数据读出来了， 并保存在了 Request.Data() 中，Handle() 只需要对数
	// 据处理即可，处理完后需要重新写入到 Request.Conn()
	//
	// 如果再次读取，会导致无限阻塞，因为 conn 的缓冲区已经读空了，参考 EchoHandler
	// 中 Handle 的注释说明
	Handle(*Request)
	AfterHandle(*Request)
}

type EchoHandler struct{}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (e *EchoHandler) BeforeHandle(r *Request) {
	//fmt.Println("before func")
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
	//fmt.Println("after func")
}

type HeartBeatHandler struct{}

func NewHeartBeatHandler() *HeartBeatHandler {
	return &HeartBeatHandler{}
}

func (h *HeartBeatHandler) BeforeHandle(req *Request) {

}

func (h *HeartBeatHandler) Handle(req *Request) {

}

func (h *HeartBeatHandler) AfterHandle(req *Request) {

}
