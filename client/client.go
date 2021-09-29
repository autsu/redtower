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
	err  error	// 将错误保存至此，以实现链式调用
}

func NewClientWithTCP(host, port string) *Client {
	c := &Client{}
	tcpaddr, _ := net.ResolveTCPAddr("tcp",
		fmt.Sprintf("%s:%s", host, port))

	conn, err := net.DialTCP("tcp", nil, tcpaddr)
	if err != nil {
		c.err = err
		return c
	}

	tcpConn := server.NewTCPConn(conn, nil)
	c.conn = tcpConn

	return c
}

func (c *Client) Init(ctx context.Context) (*Client, error) {
	go startHeartbeat(c.conn, ctx)
	return c, c.err
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

// startHeartbeat 发送心跳包给服务端
func startHeartbeat(conn server.Conn, ctx context.Context) {
	ticker := time.NewTicker(conf.SendHeartbeatTime)
	for {
		select {
		case <-ticker.C:
			heartbeat := server.NewMessage([]byte(""), server.HeartBeatMsg)
			conn.Send(heartbeat)
			log.Println("send heartbreat to server")
		case <-ctx.Done():
			log.Println("done! err: ", ctx.Err())
			return
		}
	}
}
