package server


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
