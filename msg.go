package zmsgbus

import (
	"context"
)

type Message interface {
	Ctx() context.Context
	Topic() string
	Msg() interface{}
}

// 通道消息
type channelMsg struct {
	ctx   context.Context
	topic string
	msg   interface{}
}

func newMessage(ctx context.Context, topic string, msg interface{}) Message {
	return &channelMsg{
		ctx:   ctx,
		topic: topic,
		msg:   msg,
	}
}

func (m *channelMsg) Ctx() context.Context { return m.ctx }
func (m *channelMsg) Topic() string        { return m.topic }
func (m *channelMsg) Msg() interface{}     { return m.msg }
