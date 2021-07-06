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
	Send([]byte, MessageType) (int, error)
	Receive() ([]byte, error)
}

type TCPConn struct {
	conn     *net.TCPConn
	id       uint
	IsClose  bool
	ExitChan chan bool
}

func NewTCPConn(conn *net.TCPConn, id uint) *TCPConn {
	return &TCPConn{
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
	t.conn.Close()
	t.IsClose = true
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

func (t *TCPConn) Send(data []byte, msgType MessageType) (int, error) {
	log.Printf("seng to %v ...", t.RemoteAddr().String())
	if t.IsClose {
		return 0, errors.New("read error: the connection is close")
	}

	msg := NewMessage(data, msgType)
	pack := NewDataPack()

	packmsg, err := pack.Packet(msg)
	if err != nil {
		return 0, err
	}

	n, err := t.Conn().Write(packmsg)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (t *TCPConn) Receive() ([]byte, error) {
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

	buf := make([]byte, unPack.DataLen())
	_, err = t.Conn().Read(buf)
	if err != nil {
		log.Println("read data error: ", err)
		return nil, err
	}

	return buf, nil
}
