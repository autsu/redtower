package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/autsu/redtower/conf"
	"log"
	"net"
	"sync"
	"time"

	"github.com/autsu/redtower/server"
)

type Client struct {
	once    sync.Once
	conn    server.Conn
	rawConn net.Conn
	addr    string
	err     error // 将错误保存至此，以实现链式调用
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
	c.addr = conn.LocalAddr().String()
	c.rawConn = conn

	tcpConn := server.NewTCPConn(conn, nil)
	c.conn = tcpConn
	return c
}

func (c *Client) Init(ctx context.Context) (*Client, error) {
	c.once.Do(func() {
		go startHeartbeat(c.conn, ctx)
	})
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

func (c *Client) Close() error {
	return c.rawConn.Close()
}

func (c *Client) Addr() string {
	return c.addr
}

// startHeartbeat 发送心跳包给服务端
func startHeartbeat(conn server.Conn, ctx context.Context) {
	ticker := time.NewTicker(time.Duration(conf.GlobalConf.HeartBeatSendingInterval) * time.Second)
	for {
		select {
		case <-ticker.C:
			heartbeat := server.NewMessage([]byte(""), server.HearbeatMsg)
			conn.Send(heartbeat)
			log.Println("send heartbreat to server")
		case <-ctx.Done():
			log.Println("done! err: ", ctx.Err())
			return
		}
	}
}
