package server

type MessageType int8

// 系统保留类型
const (
	ErrorMsg    MessageType = -1
	HearbeatMsg MessageType = -2
)

type Message struct {
	dataLen uint64
	data    []byte
	typ     MessageType
}

func NewMessage(data []byte, typ MessageType) *Message {
	return &Message{
		dataLen: (uint64)(len(data)),
		data:    data,
		typ:     typ,
	}
}

func NewErrorMessage(data []byte) *Message {
	return &Message{
		dataLen: (uint64)(len(data)),
		data:    data,
		typ:     ErrorMsg,
	}
}

func (m *Message) DataLen() uint64 {
	return m.dataLen
}

func (m *Message) Data() []byte {
	return m.data
}

func (m *Message) Type() MessageType {
	return m.typ
}
