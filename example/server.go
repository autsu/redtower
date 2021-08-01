package main

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"zinx/server"
)

// EchoHandler 自己实现 Echo message 的处理函数
type EchoHandler struct{}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (e *EchoHandler) BeforeHandle(r *server.Request) {}

func (e *EchoHandler) Handle(r *server.Request) {
	data := r.Data()

	msg := server.NewMessage(data, r.MsgType())
	_, err := r.Conn().Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
}

func (e *EchoHandler) AfterHandle(r *server.Request) {}

// HttpEchoPostFormHandler 回显 post 表单数据（该 handler 仅仅是用来测试 http 请求）
type HttpEchoPostFormHandler struct {

}

func NewHttpEchoPostFormHandler() *HttpEchoPostFormHandler {
	return &HttpEchoPostFormHandler{}
}

func (h *HttpEchoPostFormHandler) BeforeHandle(req *server.Request) {

}

func (h *HttpEchoPostFormHandler) Handle(req *server.Request) {
	bsr := bytes.NewReader(req.Data())
	br := bufio.NewReader(bsr)


	r, err := http.ReadRequest(br)
	if err != nil {
		log.Println(err)
		msg := server.NewMessage([]byte(err.Error()), server.ErrorMsg)
		req.Conn().Send(msg)
		return
	}
	log.Printf("%+v", r)

	r.ParseForm()
	form := r.Form.Encode()
	msg := server.NewMessage([]byte(form), server.HTTPMsg)
	req.Conn().Send(msg)
}

func (h *HttpEchoPostFormHandler) AfterHandle(req *server.Request) {

}

func init() {
	log.SetFlags(log.LstdFlags|log.Lshortfile|log.Ltime)
}

func main() {
	s := server.NewTCPServer("localhost", "8080", "server1")
	s.AddHandler(server.EchoMsg, NewEchoHandler())
	s.AddHandler(server.HeartBeat, server.NewHeartBeatHandler())
	s.AddHandler(server.HTTPMsg, NewHttpEchoPostFormHandler())

	s.Start()
}
