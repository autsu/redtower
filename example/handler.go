package example

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"zinx/server"
)

const (
	EchoMsg = iota
	HTTPMsg
)

// EchoHandler 自己实现 Echo message 的处理函数
type EchoHandler struct {
	server.BasicHandler
}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (e *EchoHandler) Handle(r *server.Request) {
	data := r.Data()

	msg := server.NewMessage(data, r.MsgType())
	_, err := r.Conn().Send(msg)
	if err != nil {
		log.Println(err)
		return
	}
}

// HttpEchoPostFormHandler 回显 post 表单数据（该 handler 仅仅是用来测试 http 请求）
type HttpEchoPostFormHandler struct {
	server.BasicHandler
}

func NewHttpEchoPostFormHandler() *HttpEchoPostFormHandler {
	return &HttpEchoPostFormHandler{}
}

func (h *HttpEchoPostFormHandler) Handle(req *server.Request) {
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
	msg := server.NewMessage([]byte(form), HTTPMsg)
	req.Conn().Send(msg)
}
