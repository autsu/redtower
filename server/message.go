package server

import (
	"crypto/md5"
	"io"
	"sync"
)

// md5 到 typeName 的映射
var md5Typ = make(map[[16]byte]string)
var typMutex sync.Mutex

type MessageType interface {
	ID() [16]byte
}

// 没有什么用处，仅用于 GenMsgTyp 的返回
type __xxx__ struct {
	// 用作 typ 的唯一标识，使用 md5 生成，因为使用 binary 传输，不允许变长类型，
	// 所以这里定义的是一个 [16]byte，这符合 md5 的固定长度。
	// typ 需要添加到 router 的 handler map 中，所以需要唯一标识
	X [16]byte
}

func (x __xxx__) ID() [16]byte {
	return x.X
}

func GenMsgTyp(typName string) MessageType {
	h := md5.New()
	io.WriteString(h, typName)
	r := h.Sum(nil)
	v := *(*[16]byte)(r)

	typMutex.Lock()
	md5Typ[v] = typName
	typMutex.Unlock()

	return __xxx__{X: v}
}

var (
	HeartBeatMsg = GenMsgTyp("heartbeat")
	ErrorMsg     = GenMsgTyp("error")
	OriginalMsg  = GenMsgTyp("original")
)

type Message struct {
	dataLen uint32
	data    []byte
	typ     MessageType
}

func NewMessage(data []byte, typ MessageType) *Message {
	return &Message{
		dataLen: (uint32)(len(data)),
		data:    data,
		typ:     typ,
	}
}

func NewErrorMessage(data []byte) *Message {
	return &Message{
		dataLen: (uint32)(len(data)),
		data:    data,
		typ:     ErrorMsg,
	}
}

func (m *Message) DataLen() uint32 {
	return m.dataLen
}

func (m *Message) Data() []byte {
	return m.data
}

func (m *Message) Type() MessageType {
	return m.typ
}

func (m *Message) SetDataLen(len uint32) {
	m.dataLen = len
}

func (m *Message) SetData(data []byte) {
	m.data = data
}

func (m *Message) SetType(t MessageType) {
	m.typ = t
}
