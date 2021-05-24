/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/3/20
   Description :
-------------------------------------------------
*/

package zmsgbus

import (
	"sync"
)

// 默认消息队列大小
const DefaultQueueSize = 1000

type MessageBus interface {
	// 发布
	Publish(topic string, msg interface{})
	// 订阅, 返回订阅号
	Subscribe(topic string, fn ProcessFn) (subscribeId uint32)
	// 取消订阅
	Unsubscribe(topic string, subscribeId uint32)
	// 关闭主题, 同时关闭所有订阅该主题的订阅者
	Close(topic string)
}

var defaultMsgBus = NewMsgBus()

// 消息总线
type msgBus struct {
	queueSize int
	topics    msgTopics
	mx        sync.RWMutex
}

func (m *msgBus) Publish(topic string, msg interface{}) {
	m.mx.RLock()
	t, ok := m.topics[topic]
	m.mx.RUnlock()

	if ok {
		t.Publish(msg)
	}
}

func (m *msgBus) Subscribe(topic string, fn ProcessFn) (subscribeId uint32) {
	m.mx.RLock()
	t, ok := m.topics[topic]
	m.mx.RUnlock()

	if !ok {
		m.mx.Lock()
		t, ok = m.topics[topic]
		if !ok {
			t = newMsgTopic(topic)
			m.topics[topic] = t
		}
		m.mx.Unlock()
	}
	return t.Subscribe(m.queueSize, fn)
}

func (m *msgBus) Unsubscribe(topic string, subscribeId uint32) {
	m.mx.RLock()
	t, ok := m.topics[topic]
	m.mx.RUnlock()

	if ok {
		t.Unsubscribe(subscribeId)
	}
}

func (m *msgBus) Close(topic string) {
	m.mx.Lock()
	t, ok := m.topics[topic]
	if ok {
		delete(m.topics, topic)
	}
	m.mx.Unlock()

	if ok {
		t.Close()
	}
}

// 创建一个消息总线
func NewMsgBus() MessageBus {
	return NewMsgBusWithQueueSize(DefaultQueueSize)
}

// 创建一个消息总线并设置队列大小
func NewMsgBusWithQueueSize(queueSize int) MessageBus {
	if queueSize < 1 {
		queueSize = DefaultQueueSize
	}

	return &msgBus{
		queueSize: queueSize,
		topics:    make(msgTopics),
	}
}
