package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
	"zinx/conf"
	"zinx/server"
)

type Client struct {
	conn server.Conn
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
	if receive.Type() == server.ErrorMsg {
		return nil, errors.New(string(receive.Data()))
	}
	return receive, nil
}

// StartHeartbeat 发送心跳包给服务端
// change 2021.8.15: 在方法内部开启 goroutine，而不是让调用方自己开启 goroutine，
// 现在只需要 c.SendHearBeat() 即可，而不是 go c.SendHearBeat()
func (c *Client) StartHeartbeat(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(conf.SendHeartbeatTime)
		for {
			select {
			case <-ticker.C:
				heartbeat := server.NewMessage([]byte(""), server.HeartBeatMsg)
				c.conn.Send(heartbeat)
				log.Println("send heartbreat to server")
			case <-ctx.Done():
				log.Println("done! err: ", ctx.Err())
				return
			}
		}
	}()
}
