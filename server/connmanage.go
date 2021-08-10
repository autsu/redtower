package server

import (
	"fmt"
	"log"
	"strings"
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
	log.Printf("cur conn manage info: %v\n", c.String())
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

func (c *ConnManage) String() string {
	var sb strings.Builder
	for id, conn := range c.conns {
		info := fmt.Sprintf("conn[id = %d] addr: %v\n", id, conn.RemoteAddr().String())
		sb.WriteString(info)
	}
	return sb.String()
}
