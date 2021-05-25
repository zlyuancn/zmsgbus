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
	Subscribe(topic string, threadCount int, handler Handler) (subscribeId uint32)
	// 全局订阅, 会收到所有消息
	SubscribeGlobal(threadCount int, handler Handler) (subscribeId uint32)
	// 取消订阅
	Unsubscribe(topic string, subscribeId uint32)
	// 取消全局订阅
	UnsubscribeGlobal(subscribeId uint32)
	// 关闭主题, 同时关闭所有订阅该主题的订阅者
	CloseTopic(topic string)
	// 关闭
	Close()
}

// 消息总线
type msgBus struct {
	global    *msgTopic // 用于接收全局消息
	queueSize int
	topics    msgTopics
	mx        sync.RWMutex // 用于锁 topics
}

func (m *msgBus) Publish(topic string, msg interface{}) {
	m.global.Publish(topic, msg) // 发送消息到全局

	m.mx.RLock()
	t, ok := m.topics[topic]
	m.mx.RUnlock()

	if ok {
		t.Publish(topic, msg)
	}
}

func (m *msgBus) Subscribe(topic string, threadCount int, handler Handler) (subscribeId uint32) {
	m.mx.RLock()
	t, ok := m.topics[topic]
	m.mx.RUnlock()

	if !ok {
		m.mx.Lock()
		t, ok = m.topics[topic]
		if !ok {
			t = newMsgTopic()
			m.topics[topic] = t
		}
		m.mx.Unlock()
	}
	return t.Subscribe(m.queueSize, threadCount, handler)
}
func (m *msgBus) SubscribeGlobal(threadCount int, handler Handler) (subscribeId uint32) {
	return m.global.Subscribe(m.queueSize, threadCount, handler)
}

func (m *msgBus) Unsubscribe(topic string, subscribeId uint32) {
	m.mx.RLock()
	t, ok := m.topics[topic]
	m.mx.RUnlock()

	if ok {
		t.Unsubscribe(subscribeId)
	}
}
func (m *msgBus) UnsubscribeGlobal(subscribeId uint32) {
	m.global.Unsubscribe(subscribeId)
}

func (m *msgBus) CloseTopic(topic string) {
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

func (m *msgBus) Close() {
	m.mx.Lock()
	for _, t := range m.topics {
		t.Close()
	}
	m.global.Close()
	m.global = newMsgTopic()
	m.topics = make(msgTopics)
	m.mx.Unlock()
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
		global:    newMsgTopic(),
		queueSize: queueSize,
		topics:    make(msgTopics),
	}
}
