package server

// Handler 自定义处理函数，只要实现该接口即可，然后通过 server.AddHandler
// 注册到路由中
type Handler interface {
	BeforeHandle(*Request)
	// Handle 不需要再从 conn 中读取了，在该方法的上层调用 conn.Handler() 中
	// 已经将数据读出来了， 并保存在了 Request.Data() 中，Handle() 只需要对数
	// 据处理即可，处理完后需要重新写入到 Request.Conn()
	//
	// 如果再次读取，会导致无限阻塞，因为 conn 的缓冲区已经读空了，参考 EchoHandler
	// 中 Handle 的注释说明
	Handle(*Request)
	AfterHandle(*Request)
}

type HeartBeatHandler struct{}

func NewHeartBeatHandler() *HeartBeatHandler {
	return &HeartBeatHandler{}
}

func (h *HeartBeatHandler) BeforeHandle(req *Request) {

}

func (h *HeartBeatHandler) Handle(req *Request) {

}

func (h *HeartBeatHandler) AfterHandle(req *Request) {

}
