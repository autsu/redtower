package server

// Pool 一个队列对应一个 goroutine，goroutine 不断从队列中取出任务并执行
type Pool struct {
	Size      int             // 最大的 goroutine 数量
	WorkQueue []chan *Request // 任务队列，队列中存储多个请求
}

func NewPool(size int) *Pool {
	return &Pool{
		Size:      size,
		WorkQueue: make([]chan *Request, size),
	}
}

func (p *Pool) AddWork(r *Request) {

}
