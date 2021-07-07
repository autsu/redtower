package server

import (
	"errors"
	"io"
	"log"
	"net"
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
	server   *TCPServer
	conn     *net.TCPConn
	id       uint
	IsClose  bool
	ExitChan chan bool
}

func NewTCPConn(conn *net.TCPConn, server *TCPServer, id uint) *TCPConn {
	return &TCPConn{
		server:   server,
		conn:     conn,
		id:       id,
		IsClose:  false,
		ExitChan: make(chan bool),
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

	if err := t.server.Router.Do(req); err != nil {
		log.Println("router do func error: ", err)
		return err
	}

	return nil
}
