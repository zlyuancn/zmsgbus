/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/3/20
   Description :
-------------------------------------------------
*/

package zmsgbus

import (
	"context"
	"sync"
)

// 默认消息队列大小
const DefaultMsgQueueSize = 1000

type MessageBus interface {
	// 发布
	Publish(ctx context.Context, topic string, msg interface{})
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
	global       Topic // 用于接收全局消息
	msgQueueSize int
	topics       map[string]Topic
	mx           sync.RWMutex // 用于锁 topics
}

// 创建一个消息总线
func NewMsgBus() MessageBus {
	return NewMsgBusWithQueueSize(DefaultMsgQueueSize)
}

// 创建一个消息总线并设置消息队列缓存大小, 消息队列满时会阻塞消息发送
func NewMsgBusWithQueueSize(msgQueueSize int) MessageBus {
	if msgQueueSize < 1 {
		msgQueueSize = DefaultMsgQueueSize
	}

	return &msgBus{
		global:       newMsgTopic(),
		msgQueueSize: msgQueueSize,
		topics:       make(map[string]Topic),
	}
}

func (m *msgBus) Publish(ctx context.Context, topic string, msg interface{}) {
	m.global.Publish(ctx, topic, msg) // 发送消息到全局

	m.mx.RLock()
	t, ok := m.topics[topic]
	m.mx.RUnlock()

	if ok {
		t.Publish(ctx, topic, msg)
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
	return t.Subscribe(m.msgQueueSize, threadCount, handler)
}
func (m *msgBus) SubscribeGlobal(threadCount int, handler Handler) (subscribeId uint32) {
	return m.global.Subscribe(m.msgQueueSize, threadCount, handler)
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

	// 后置关闭
	if ok {
		t.Close()
	}
}

func (m *msgBus) Close() {
	m.mx.Lock()
	clearTopic := make([]Topic, 0, 1+len(m.topics))
	clearTopic = append(clearTopic, m.global)
	for _, t := range m.topics {
		clearTopic = append(clearTopic, t)
	}
	m.global = newMsgTopic()
	m.topics = make(map[string]Topic)
	m.mx.Unlock()

	// 后置关闭
	for _, t := range clearTopic {
		t.Close()
	}
}
