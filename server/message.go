package server

type MessageType int32

const (
	EchoMsg MessageType = iota
	TextMsg
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

type Message interface {
	DataLen() uint32
	Data() []byte
	Type() MessageType

	SetDataLen(uint32)
	SetData([]byte)
	SetType(MessageType)
}

type message struct {
	dataLen uint32
	data    []byte
	type_   MessageType
}

func NewMessage(data []byte, type_ MessageType) *message {
	return &message{
		dataLen: (uint32)(len(data)),
		data:    data,
		type_:   type_,
	}
}

func (m *message) DataLen() uint32 {
	return m.dataLen
}

func (m *message) Data() []byte {
	return m.data
}

func (m *message) Type() MessageType {
	return m.type_
}

func (m *message) SetDataLen(len uint32) {
	m.dataLen = len
}

func (m *message) SetData(data []byte) {
	m.data = data
}

func (m *message) SetType(t MessageType) {
	m.type_ = t
}
