package server

import (
	"log"
	"sync"
)

type ConnManage struct {
	conns map[uint64]Conn
	mu sync.RWMutex
}

func NewConnManage() *ConnManage {
	return &ConnManage{
		conns: make(map[uint64]Conn),
		mu:    sync.RWMutex{},
	}
}

func (c *ConnManage) Add(conn Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conns[conn.ConnID()] = conn

	log.Printf("conn[id = %d, addr = %s] add to manage\n",
		conn.ConnID(), conn.RemoteAddr().String())
}

func (c *ConnManage) Remove(conn Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.conns, conn.ConnID())

	log.Printf("conn[id = %v, addr = %v] remove from manage\n",
		conn.ConnID(), conn.RemoteAddr())
}

func (c *ConnManage) GetConnById(connId uint64) Conn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conns[connId]
}

func (c *ConnManage) All() map[uint64]Conn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conns
}

func (c *ConnManage) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.conns)
}
