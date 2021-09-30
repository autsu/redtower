package server

import (
	"log"
	"time"
	"github.com/zengh1/redtower/conf"
)

type HeartBeat struct {
	conn Conn
}

func NewHeartBeat(conn Conn) *HeartBeat {
	return &HeartBeat{
		conn: conn,
	}
}

func (h *HeartBeat) Start() {
	log.Println("🫀 start Heartbeat detection...")
	for {
		select {
		case <-time.After(conf.DefaultDeadlineTime):
			h.conn.Stop()
			log.Printf("⚠️ haven't received the heartbeat packet for a long time,"+
				" the conn[id = %v] is close \n", h.conn.ConnID())
			return
		// conn.Handler() 中如果读取到了数据，会发送信号到 t.HeartbeatChan
		// 这里就可以读取出来了
		case <-h.conn.HeartBeatChan():
			log.Println("receive heartbeat from client")
		}
	}

}
