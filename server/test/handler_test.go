package test

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"
	"zinx/server"
)

func TestEchoHandler(t *testing.T) {
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

		tcpConn := server.NewTCPConn(conn, 1)

		request := server.NewRequest(server.NewMessage(nil, 1), tcpConn)

		echoHandler := server.NewEchoHandler()

		echoHandler.Handle(request)
	}
}

func TestEchoHandlerDial(t *testing.T) {
	tcpaddr, _ := net.ResolveTCPAddr("tcp", ":8080")
	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		log.Fatalln(err)
	}

	tcpConn := server.NewTCPConn(conn, 1)
	msgs := []string{"666", "456", "789"}

	for {
		// scan 只能扫描到 main goroutine，test 中无法扫描到
		//var input string
		//_, err = fmt.Scan(&input)
		//if err != nil {
		//	log.Println(err)
		//	break
		//}

		for _, msg := range msgs {
			n, err := tcpConn.Send([]byte(msg), server.EchoMsg)
			log.Printf("send %d bytes", n)
			if err != nil {
				log.Fatalln(err)
			}

			receive, err := tcpConn.Receive()
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(string(receive))
			time.Sleep(time.Second * 3)
		}

	}

}
