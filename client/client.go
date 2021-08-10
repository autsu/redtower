package client

import (
	"fmt"
	"log"
	"net"
	"time"
	"zinx/conf"
	"zinx/server"
)

type Client struct {
	conn     server.Conn
}

func NewClientWithTCP(dialHost, dialPort string) *Client {
	tcpaddr, _ := net.ResolveTCPAddr("tcp",
		fmt.Sprintf("%s:%s", dialHost, dialPort))

	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		log.Fatalln(err)
	}

	tcpConn := server.NewTCPConn(conn, nil)

	return &Client{
		conn: tcpConn,
	}
}

func (c *Client) Send(msg *server.Message) (int, error) {
	n, err := c.conn.Send(msg)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (c *Client) Receive() (*server.Message, error) {
	receive, err := c.conn.Receive()
	if err != nil {
		return nil, err
	}
	return receive, nil
}

func (c *Client) SendHeartbeat() {
	ticker := time.NewTicker(conf.SendHeartbeatTime)
	for {
		select {
		case <-ticker.C:
			heartbeat := server.NewMessage([]byte(""), server.HeartBeatMsg)
			c.conn.Send(heartbeat)
			//log.Println("send heartbreat to server")
		}
	}
}
