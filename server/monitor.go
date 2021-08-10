package server

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"time"
)

type Monitor struct {
	server Server
}

func NewMonitor(server Server) *Monitor {
	return &Monitor{server: server}
}

// Start 每隔 interval 时间，就记录当前状态到 output 中
func (m *Monitor) Start(output io.Writer, interval time.Duration) {
	s := status(m.server)
	t := time.NewTicker(interval)
	for {
		select {
		case <-t.C :
			io.Copy(output, strings.NewReader(s))
		}
	}
}

func status(server Server) string {
	var sb strings.Builder

	sb.WriteString("======== 监控信息 ==========\n")
	s := fmt.Sprintf("连接数：%d \n", server.ConnManage().Len())
	sb.WriteString(s)
	sb.WriteString(server.ConnManage().String())
	sb.WriteString("\n")
	s = fmt.Sprintf("goroutine 数量：%d\n", runtime.NumGoroutine())
	sb.WriteString(s)

	sb.WriteString("===========================\n")

	return sb.String()
}


