package server

import (
	"github.com/autsu/redtower/logs"
	"sync"
)

var GlobalConnManage = newConnManage()

type ConnManage struct {
	// k: string，存储 server
	// v: map，存储某个 server 下的所有 connection
	conns map[string]map[uint64]Conn
	mu    sync.RWMutex
}

func newConnManage() *ConnManage {
	return &ConnManage{
		conns: make(map[string]map[uint64]Conn),
	}
}

func (c *ConnManage) Add(serverName string, conn Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.conns[serverName]
	if !ok {
		c.conns[serverName] = make(map[uint64]Conn)
	}
	c.conns[serverName][conn.ConnID()] = conn

	logs.L.Printf_IfOpenDebug("conn[id=%v, addr=%v] add to manage[server=%v]\n",
		conn.ConnID(), conn.Addr().String(), serverName)
}

func (c *ConnManage) RemoveAndClose(serverName string, conn Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.conns[serverName], conn.ConnID())
	if conn != nil {
		conn.Close()
		logs.L.Printf_IfOpenDebug("conn[id=%v, addr=%v] remove from manage[server=%v]\n",
			conn.ConnID(), conn.Addr().String(), serverName)
	}
}

func (c *ConnManage) GetConn(serverName string, connId uint64) (conn Conn, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	conn, ok = c.conns[serverName][connId]
	return
}

func (c *ConnManage) GetServerAllConn(serverName string) map[uint64]Conn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conns[serverName]
}

func (c *ConnManage) Nums() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return (uint64)(len(c.conns))
}

func (c *ConnManage) NumsFromServer(serverName string) uint64 {
	return (uint64)(len(c.conns[serverName]))
}

func (c *ConnManage) ClearAll() {
	for _, conn := range c.conns {
		for _, c := range conn {
			c.Close()
		}
	}
}

func (c *ConnManage) ClearFromServer(serverName string) {
	for _, conn := range c.conns[serverName] {
		conn.Close()
	}
}
