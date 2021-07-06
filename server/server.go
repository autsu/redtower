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
	Port string
	Host string
	Name string
}

func NewTCPServer(host, port, name string) *TCPServer {
	return &TCPServer{
		Port: port,
		Host: host,
		Name: name,
	}
}

func (t *TCPServer) Start() {
	log.Println("server start...")
	addr := t.Host + ":" + t.Port
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("listen socket error: ", err)
		return
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("accept conn error: ", err)
			continue
		}

		go func() {
			buf := make([]byte, 512)

			for {
				n, err := conn.Read(buf)
				if err != nil {
					log.Println("read from conn error: ", err)
					return
				}
				log.Printf("read %d bytes\n", n)

				_, err = conn.Write(buf[:n])
				if err != nil {
					log.Println("write to conn error: ", err)
					return
				}
			}

		}()

	}


}

func (t *TCPServer) Stop() {

}

func (t *TCPServer) Server() {

}
