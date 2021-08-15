package main

import (
	"bufio"
	"bytes"
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"zinx/server"
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

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Ltime)
}

func main() {
	// pprof
	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	s := server.NewTCPServer("localhost", "8080", "server1")
	s.AddHandler(server.EchoMsg, NewEchoHandler())
	s.AddHandler(server.HeartBeatMsg, server.NewHeartBeatHandler())
	s.AddHandler(server.HTTPMsg, NewHttpEchoPostFormHandler())

	//go func() {
	//	// 开启监控
	//	monitor := server.NewMonitor(s)
	//	monitor.Start(os.Stdout, time.Second * 10)
	//}()

	s.Start(context.Background())
}
