package zmsgbus

import (
	"sync/atomic"
)

// 全局自增订阅者id
var autoIncrSubscriberId uint32

// 生成下一个订阅者id
func nextSubscriberId() uint32 {
	return atomic.AddUint32(&autoIncrSubscriberId, 1)
}

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
