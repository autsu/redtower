package server

import (
	"log"
	"sync"
)

type ConnManage struct {
	conns map[uint64]Conn
	Mu sync.RWMutex
}

func NewConnManage() *ConnManage {
	return &ConnManage{
		conns: make(map[uint64]Conn),
		Mu:    sync.RWMutex{},
	}
}

func (c *ConnManage) Add(conn Conn) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.conns[conn.ConnID()] = conn

	log.Printf("conn[id = %d, addr = %s] add to manage\n",
		conn.ConnID(), conn.RemoteAddr().String())
}

func (c *ConnManage) Remove(conn Conn) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	delete(c.conns, conn.ConnID())

	log.Printf("conn[id = %d, addr = %d] remove from manage\n",
		conn.ConnID(), conn.RemoteAddr())
}

func (c *ConnManage) GetConnById(connId uint64) Conn {
	return c.conns[connId]
}

func (c *ConnManage) Len() int {
	return len(c.conns)
}
