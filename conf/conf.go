package conf

import "time"

const (
	WaitHeartbeatTime   = time.Second * 10 // 等待客户端响应心跳包时间
	DefaultDeadlineTime = time.Second * 5 // 默认触发心跳检测时间
)
