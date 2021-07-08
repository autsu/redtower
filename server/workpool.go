package server

import (
	"math/rand"
	"time"
)

// DefaultQueueSize 队列的默认任务数
const DefaultQueueSize = 10

// Pool 一个队列对应一个 goroutine，goroutine 不断从队列中取出任务并执行
type Pool struct {
	Size      int64             // 最大的 goroutine 数量
	WorkQueue []chan *Request // 任务队列，队列中存储多个请求
}

func NewPool(size int64) *Pool {
	return &Pool{
		Size:      size,
		WorkQueue: make([]chan *Request, DefaultQueueSize),
	}
}

func (p *Pool) AddWork(r *Request) {
	rand.Seed(time.Now().Unix())

	n := rand.Int63n(p.Size)

	p.WorkQueue[n] <- r
}
