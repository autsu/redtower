package server

import (
	"log"
	"time"
	"zinx/conf"
)


func (t *TCPConn) Heartbeat() {
	log.Println("start Heartbeat...")
	for {
		select {
		case <-time.After(conf.DefaultDeadlineTime):
			log.Println("Haven't received the heartbeat packet for a long time, the conn is close")
			t.Stop()
			return
		// conn.Handler() 中如果读取到了数据，会发送信号到 t.HeartbeatChan
		// 这里就可以读取出来了
		case <-t.HeartbeatChan:
			log.Println("receive heartbeat from client")
		}
	}

}
