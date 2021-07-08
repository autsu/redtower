package server

import (
	"errors"
	"io"
	"log"
	"net"
	"time"
	"zinx/conf"
)

var (
	HeartbeatRequestChan  = make(chan struct{}, 1) // 服务端通知客户端发送心跳
	HeartbeatResponseChan = make(chan struct{}, 1) // 客户端的心跳响应
)

type Conn interface {
	Start()
	Stop()
	Conn() net.Conn
	RemoteAddr() net.Addr
	ConnID() uint
	Send(data *Message) (int, error)
	Receive() (*Message, error)
}

type TCPConn struct {
	server  *TCPServer
	conn    *net.TCPConn
	id      uint
	IsClose bool
	// 连接最后的活动时间，该字段主要用于连接的保活
	// 每当进行 io 操作后更新该值为 time.Now()，同时开启一条线程不断监听该字段，
	// 如果 LastActivityTime + outTime 时间内没有收到数据，则发送心跳包，
	// 如果心跳包无响应则代表对端已断开，此时可以关闭连接
	// 该字段避免了持续性的发送心跳包
	LastActivityTime time.Time
}

func NewTCPConn(conn *net.TCPConn, server *TCPServer, id uint) *TCPConn {
	return &TCPConn{
		server:           server,
		conn:             conn,
		id:               id,
		IsClose:          false,
		LastActivityTime: time.Now(),
	}
}

func (t *TCPConn) Start() {

}

func (t *TCPConn) Stop() {
	if t.IsClose {
		return
	}
	t.IsClose = true
	t.conn.Close()

}

func (t *TCPConn) Conn() net.Conn {
	return t.conn
}

func (t *TCPConn) RemoteAddr() net.Addr {
	return t.conn.RemoteAddr()
}

func (t *TCPConn) ConnID() uint {
	return t.id
}

func (t *TCPConn) Send(data *Message) (int, error) {
	log.Printf("seng to %v ...", t.RemoteAddr().String())
	if t.IsClose {
		return 0, errors.New("read error: the connection is close")
	}

	pack := NewDataPack()

	packmsg, err := pack.Packet(data)
	if err != nil {
		return 0, err
	}

	n, err := t.Conn().Write(packmsg)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (t *TCPConn) Receive() (*Message, error) {
	log.Printf("receive from %v ...", t.RemoteAddr().String())
	if t.IsClose {
		return nil, errors.New("read error: the connection is close")
	}

	pack := NewDataPack()

	headBuf := make([]byte, pack.HeadSize())

	_, err := io.ReadFull(t.Conn(), headBuf)
	if err != nil {
		log.Println("read head error: ", err)
		return nil, err
	}

	unPack, err := pack.UnPack(headBuf)
	if err != nil {
		log.Println("unpack error: ", err)
		return nil, err
	}

	dataLen := unPack.DataLen()
	msgType := unPack.Type()

	buf := make([]byte, dataLen)
	_, err = t.Conn().Read(buf)
	if err != nil {
		log.Println("read data error: ", err)
		return nil, err
	}

	m := NewMessage(buf, msgType)
	return m, nil
}

func (t *TCPConn) Handler(msg *Message) error {
	req := NewRequest(msg, t)

	go t.Keepalive()

	if err := t.server.Router.Do(req); err != nil {
		log.Println("router do func error: ", err)
		return err
	}

	return nil
}

func (t *TCPConn) Keepalive() {
	for {
		// 如果最后活动时间 + 超时时间 小于 当前时间，则发送心跳包检测对方是否存活
		if t.LastActivityTime.Add(conf.DefaultDeadlineTime).Before(time.Now()) {
			log.Println("timeout, send heartbeat to client")
			HeartbeatRequestChan <- struct{}{}
			log.Println("send request to chan")

			for {
				select {
				// 等待 WaitHeartbeatTime 后还没有收到心跳包，则关闭连接
				case <-time.After(conf.WaitHeartbeatTime):
					log.Println("tcp conn timeout")
					t.Conn().Close()
					t.IsClose = true
					return
				// 收到了客户端的响应，更新最后活动时间
				case <-HeartbeatResponseChan:
					log.Println("Receive a heartbeat packet from the client")
					t.LastActivityTime = time.Now()
					return
				default:

				}
			}
		}
		time.Sleep(time.Second)
	}
}
