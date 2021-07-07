package server

type MessageType int32

const (
	EchoMsg MessageType = iota
	TextMsg
	HTTPMsg
)

func TypeOfMessage(msgType MessageType) string {
	switch msgType {
	case EchoMsg:
		return "回响信息"
	case TextMsg:
		return "文本信息"
	}
	return "未知种类的信息"
}

//type Message interface {
//	DataLen() uint32
//	Data() []byte
//	Type() MessageType
//
//	SetDataLen(uint32)
//	SetData([]byte)
//	SetType(MessageType)
//}

type Message struct {
	dataLen uint32
	data    []byte
	type_   MessageType
}

func NewMessage(data []byte, type_ MessageType) *Message {
	return &Message{
		dataLen: (uint32)(len(data)),
		data:    data,
		type_:   type_,
	}
}

func (m *Message) DataLen() uint32 {
	return m.dataLen
}

func (m *Message) Data() []byte {
	return m.data
}

func (m *Message) Type() MessageType {
	return m.type_
}

func (m *Message) SetDataLen(len uint32) {
	m.dataLen = len
}

func (m *Message) SetData(data []byte) {
	m.data = data
}

func (m *Message) SetType(t MessageType) {
	m.type_ = t
}
