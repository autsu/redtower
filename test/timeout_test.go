package test

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"
	"zinx/server"
)

func TestTimeoutServer(t *testing.T) {
	s := server.NewTCPServer("localhost", "8080", "server1")
	s.AddHandler(server.EchoMsg, server.NewEchoHandler())
	s.AddHandler(server.HeartBeat, server.NewHeartBeatHandler())

	s.Start()
}

func TestTimeOutDial(t *testing.T) {
	tcpaddr, _ := net.ResolveTCPAddr("tcp", ":8080")
	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		log.Fatalln(err)
	}

	tcpConn := server.NewTCPConn(conn, nil, 1)
	msgs := []string{"666", "456", "789"}

	//go func() {
	//	for {
	//		select {
	//		case <-server.HeartbeatRequestChan:
	//			log.Println("client receive server heartbeat request")
	//			server.HeartbeatResponseChan <- struct{}{}
	//		default:
	//
	//		}
	//	}
	//}()

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
		time.Sleep(time.Second * 15)
	}

}
