package test

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"
	"zinx/server"
)

// 完全没有封装
func TestEchoHandlerServer(t *testing.T) {
	tcpaddr, _ := net.ResolveTCPAddr("tcp", ":8080")
	listen, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		s := server.NewTCPServer("localhost", "8080", "server1")
		s.AddHandler(server.EchoMsg, server.NewEchoHandler())

		tcpConn := server.NewTCPConn(conn, s, 1)

		request := server.NewRequest(server.NewMessage(nil, 1), tcpConn)

		echoHandler := server.NewEchoHandler()

		go echoHandler.Handle(request)
	}
}

// 初步封装
func TestEchoHandlerServerV2(t *testing.T) {
	tcpaddr, _ := net.ResolveTCPAddr("tcp", ":8080")
	listen, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}

		s := server.NewTCPServer("localhost", "8080", "server1")
		s.AddHandler(server.EchoMsg, server.NewEchoHandler())
		//log.Println(s.Router)

		tcpConn := server.NewTCPConn(conn, s, 1)
		for {
			// 废弃
			if err := tcpConn.Handler(); err != nil {
				log.Println(err)
				break
			}
		}

	}
}

// 完全封装
func TestEchoHandlerServerV3(t *testing.T) {
	s := server.NewTCPServer("localhost", "8080", "server1")
	s.AddHandler(server.EchoMsg, server.NewEchoHandler())

	s.Start()
}

func TestEchoHandlerDial(t *testing.T) {
	tcpaddr, _ := net.ResolveTCPAddr("tcp", ":8080")
	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		log.Fatalln(err)
	}


	tcpConn := server.NewTCPConn(conn, nil, 1)
	msgs := []string{"666", "456", "789"}

	for {
		// scan 只能扫描到 main goroutine，fff 中无法扫描到
		//var input string
		//_, err = fmt.Scan(&input)
		//if err != nil {
		//	log.Println(err)
		//	break
		//}

		for _, msg := range msgs {
			msg := server.NewMessage([]byte(msg), server.EchoMsg)
			n, err := tcpConn.Send(msg)
			log.Printf("send %d bytes", n)
			if err != nil {
				log.Fatalln(err)
			}

			receive, err := tcpConn.Receive()
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(string(receive.Data()))
			time.Sleep(time.Second * 3)
		}

	}

}

func TestEchoHandlerDial1(t *testing.T) {
	tcpaddr, _ := net.ResolveTCPAddr("tcp", ":8080")
	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		log.Fatalln(err)
	}


	tcpConn := server.NewTCPConn(conn, nil, 1)
	msgs := []string{"666", "456", "789"}

	for {
		// scan 只能扫描到 main goroutine，fff 中无法扫描到
		//var input string
		//_, err = fmt.Scan(&input)
		//if err != nil {
		//	log.Println(err)
		//	break
		//}

		for _, msg := range msgs {
			msg := server.NewMessage([]byte(msg), server.EchoMsg)
			n, err := tcpConn.Send(msg)
			log.Printf("send %d bytes", n)
			if err != nil {
				log.Fatalln(err)
			}

			receive, err := tcpConn.Receive()
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(string(receive.Data()))
			time.Sleep(time.Second * 3)
		}

	}

}
