package conf

import "time"

const (
	DefaultDeadlineTime    = time.Second * 30 // 默认触发心跳检测时间
	SendHeartbeatTime      = time.Second * 5  // 多久发送一次心跳
	MaxConnNum             = 1000             // 最大连接数量
	DefaultQueueSize       = 10               // 队列的默认任务数
	DefaultGoroutineMaxNum = 500              // workpool 默认启动的最大 goroutine 数量
)
