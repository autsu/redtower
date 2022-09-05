package example

// 自定义消息类型，需要保证不与库的内置消息类型（ErrorMsg = -1 和 HearbeatMsg = -2 ）的值相同
const (
	HTTPMsg = 10
	EchoMsg = 11
)
