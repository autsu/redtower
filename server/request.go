package server

type Request interface {
	Conn() Conn
	Data() []byte
	MsgType() MessageType
}

type request struct {
	msg  Message
	conn Conn
}

func NewRequest(msg Message, conn Conn) *request {
	return &request{
		msg:  msg,
		conn: conn,
	}
}

func (r *request) Conn() Conn {
	return r.conn
}

func (r *request) Data() []byte {
	return r.msg.Data()
}

func (r *request) MsgType() MessageType {
	return r.msg.Type()
}
