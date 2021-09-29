package server

import "math"

type MessageType int32

const (
	HeartBeatMsg MessageType = iota + math.MinInt32	// 防止和用户自定义类型冲突
	ErrorMsg
	OriginalMsg // 不进行任何处理
)

//func TypeOfMessage(msgType MessageType) string {
//	switch msgType {
//	case ErrorMsg:
//		return "错误信息"
//	case HeartBeatMsg:
//		return "心跳信息"
//	case OriginalMsg:
//		return "原始信息"
//	}
//	return "未知种类的信息"
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

func NewErrorMessage(data []byte) *Message {
	return &Message{
		dataLen: (uint32)(len(data)),
		data:    data,
		type_:   ErrorMsg,
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
