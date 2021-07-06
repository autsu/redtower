package server

import (
	"log"
	"net"
)

type Server interface {
	Start()
	Stop()
	Server()
}

type TCPServer struct {
	Port   string
	Host   string
	Name   string
	Router *Router
}

func NewTCPServer(host, port, name string) *TCPServer {
	return &TCPServer{
		Port:   port,
		Host:   host,
		Name:   name,
		Router: NewRouter(),
	}
}

func (t *TCPServer) Start() {
	log.Println("server start...")

	addr := t.Host + ":" + t.Port

	tcpaddr, _ := net.ResolveTCPAddr("tcp", addr)

	listen, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		log.Println("listen socket error: ", err)
		return
	}

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println("accept conn error: ", err)
			continue
		}

		tcpConn := NewTCPConn(conn, t, 1)
		for {
			if err := tcpConn.ReceiveAndHandler(); err != nil {
				log.Println(err)
				break
			}
		}

	}

}

func (t *TCPServer) Stop() {

}

func (t *TCPServer) Server() {

}

func (t *TCPServer) AddHandler(typ MessageType, handler Handler) {
	t.Router.AddRouter(typ, handler)
}
