package zmsgbus

// 通道消息
type channelMsg struct {
	topic string
	msg   interface{}
}

// 处理函数
type Handler func(topic string, msg interface{})

// 订阅者
type subscriber struct {
	handler Handler
	queue   chan *channelMsg
}
