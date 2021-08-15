package server

import (
	"log"
	"time"
	"zinx/conf"
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
	log.Println("start Heartbeat...")
	for {
		select {
		case <-time.After(conf.DefaultDeadlineTime):
			log.Printf("Haven't received the heartbeat packet for a long time,"+
				" the conn[id = %v] is close \n", h.conn.ConnID())
			h.conn.Stop()
			return
		// conn.Handler() 中如果读取到了数据，会发送信号到 t.HeartbeatChan
		// 这里就可以读取出来了
		case <-h.conn.HeartBeatChan():
			log.Println("receive heartbeat from client")
			// BUG: 加了 default 导致 CPU 占用达到 300-400%
			// 原因：因为有了 default 会导致 select 不会被阻塞，从而外层死循环不断执行，
			// 造成 CPU 空转
			//default:
			//	// 如果连接已经关闭则停止发送
			//	if h.conn.IsClose() {
			//		log.Printf("conn[id = %v, add = %v] is closed, stop heartbeat\n",
			//			h.conn.ConnID(), h.conn.RemoteAddr())
			//		return
			//	}
		}
	}

}
