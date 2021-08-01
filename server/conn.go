package server

import (
	"errors"
	"io"
	"log"
	"net"
	"sync/atomic"
)

type Conn interface {
	Stop()
	Conn() net.Conn
	RemoteAddr() net.Addr
	ConnID() uint64
	Send(data *Message) (int, error) // 发送数据到 conn
	Receive() (*Message, error)      // 从 conn 中接收数据
	Handler() error                  // 处理 conn 中的数据
}

// 每创建一条 conn，便从这里取得一个新的 connId，之后 connId++（原子的）
var connId uint64 = 1

type TCPConn struct {
	server        *TCPServer
	conn          *net.TCPConn
	id            uint64
	IsClose       bool
	HeartbeatChan chan struct{}
}

func NewTCPConn(conn *net.TCPConn, server *TCPServer) *TCPConn {
	t :=  &TCPConn{
		server:        server,
		conn:          conn,
		id:            atomic.LoadUint64(&connId),
		IsClose:       false,
		HeartbeatChan: make(chan struct{}),
	}
	// global conn id + 1
	for !atomic.CompareAndSwapUint64(&connId, connId, connId+1) {

	}

	return t
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

func (t *TCPConn) ConnID() uint64 {
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

func (t *TCPConn) Handler() error {
	recvData, err := t.Receive()
	if err != nil {
		log.Println(err)
		return err
	}
	// 接收到了数据，则发送信号到 t.HeartbeatChan
	t.HeartbeatChan <- struct{}{}

	msgType := recvData.Type()
	log.Printf("recvData: data: %v, type: %v",
		string(recvData.Data()), TypeOfMessage(msgType))

	msg := NewMessage(recvData.Data(), msgType)
	req := NewRequest(msg, t)

	t.server.Pool.AddWork(req)

	return nil
}
