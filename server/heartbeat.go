package server

import (
	"github.com/autsu/redtower/conf"
	"github.com/autsu/redtower/logs"
	"log"
	"time"
)

type HeartBeat struct {
	conn       Conn
	serverName string
}

func NewHeartBeat(serverName string, conn Conn) *HeartBeat {
	return &HeartBeat{
		conn:       conn,
		serverName: serverName,
	}
}

// 当 read 到数据时调用 trigger，说明这条连接是活跃的，向这条连接的心跳 chan 写入数据用作通知
func (h *HeartBeat) trigger() {
	tc, ok := h.conn.(*TCPConn)
	if !ok {
		log.Println("Currently only support tcp")
		return
	}
	tc.heartbeatChan <- struct{}{}
}

func (h *HeartBeat) Start() {
	log.Println("🫀 start Heartbeat detection...")
	tc, ok := h.conn.(*TCPConn)
	if !ok {
		log.Println("Currently only support tcp")
		return
	}

	for {
		select {
		case <-time.After(time.Duration(conf.GlobalConf.HeartBeatDeadline) * time.Second):
			GlobalConnManage.RemoveAndClose(h.serverName, h.conn)
			logs.L.Printf_IfOpenDebug(`
				⚠️ haven't received the heartbeat packet for a long time,
				the conn[id = %v] is close \n`,
				h.conn.ConnID())
			return
		// 如果从连接中读取到了数据，会发送信号到 t.HeartbeatChan
		// 从而可以走到这个分支
		case <-tc.heartbeatChan:
			logs.L.Printf_IfOpenDebug("receive heartbeat from client\n")
		}
	}
}
