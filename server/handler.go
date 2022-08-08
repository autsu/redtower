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

// HeartBeatHandler 心跳包的处理函数
type heartBeatHandler struct {
	BasicHandler
}

func (h *heartBeatHandler) Handle(req *Request) {
	// 接收到了心跳包，则发送信号到 conn.HeartbeatChan
	req.Conn().HeartBeatChan() <- struct{}{}
}

type errorHandler struct {
	BasicHandler
}

type originalHandler struct {
	BasicHandler
}

// BasicHandler 空实现，可以通过组合该结构体来实现接口，
// 这样可以让某些情况下实现接口更简洁，比如：
//
// 一个需要实现 Handler 的结构体，但是只需要实现 Handler.Handle()，其他两个方法为空，
// 此时只需要这么写即可：
//
// type A struct {
//		BasicHandler
// }
//
// 重写 Handler 方法即可
// func (a *A) Handler(r *Request) {
// 		// .... 实现内容
// }
//
// 如果没有 BasicHandler，则需要实现三个方法，即便 BeforeHandle 和 AfterHandle 没有
// 任何内容，写起来繁琐冗余，如果直接内嵌 Handler 接口会因为没有实现部分方法而导致空指针
type BasicHandler struct{}

func (b *BasicHandler) BeforeHandle(req *Request) {}

func (b *BasicHandler) Handle(req *Request) {}

func (b *BasicHandler) AfterHandle(req *Request) {}
