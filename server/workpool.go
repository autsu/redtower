package server

import "zinx/conf"

// +-----+-----+-----+
// +  G  +  G  +  G  +		goroutine
// +-----+-----+-----+
//    |		|	  |			每个 G 对应一个 Q，从 Q 中不断取出任务并执行
//    |     |     |
// +-----+-----+-----+
// +  Q  +  Q  +  Q  +		queue
// +-----+-----+-----+
// + req + req + req +
// + req + req + req +
// + req + req + req +
//

// Pool 一个队列对应一个 goroutine，goroutine 不断从队列中取出任务并执行
type Pool struct {
	Size      uint64          // 最大的 goroutine 数量
	WorkQueue []chan *Request // 任务队列，队列中存储多个请求
	R         *Router
}

func NewPool(size uint64, r *Router) *Pool {
	return &Pool{
		Size:      size,
		WorkQueue: make([]chan *Request, size),
		R:         r,
	}
}

// AddWork 添加任务到任务队列，这里使用了轮询负载均衡，保证每个任务队列的任务数平均
func (p *Pool) AddWork(r *Request) {
	index := r.Conn().ConnID() % p.Size

	p.WorkQueue[index] <- r
}

func (p *Pool) StartWorkerPool() (err error) {
	for i := 0; i < int(p.Size); i++ {
		p.WorkQueue[i] = make(chan *Request, conf.DefaultQueueSize)
		go func(i int) {
			err = p.doWork(p.WorkQueue[i])
		}(i)
	}
	return
}

func (p *Pool) doWork(c chan *Request) error {
	for {
		select {
		case req := <-c:
			if err := p.R.Do(req); err != nil {
				return err
			}
		}
	}
}
