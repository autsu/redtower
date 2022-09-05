package server

// Request 包装了 conn 和 msg，用于 onRequest 回调
// 比如一个 echo 程序，它的 onRequest 回调需要先从 msg 中读取消息，然后再通过 conn
// 发送给对方，也就是说同时需要 Message 和 Conn 两样东西，所以需要封装成一个 Request
type Request struct {
	msg  *Message
	conn Conn
}

func NewRequest(msg *Message, conn Conn) *Request {
	return &Request{
		msg:  msg,
		conn: conn,
	}
}

func (r *Request) Conn() Conn {
	return r.conn
}

func (r *Request) Data() []byte {
	return r.msg.Data()
}

func (r *Request) MsgType() MessageType {
	return r.msg.Type()
}
