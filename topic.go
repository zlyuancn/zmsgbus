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
	"sync/atomic"
)

// 主题们
type msgTopics map[string]*msgTopic

// 主题
type msgTopic struct {
	subId uint32
	subs  map[uint32]*subscriber
	mx    sync.RWMutex // 用于锁 subs
}

func newMsgTopic() *msgTopic {
	return &msgTopic{
		subs: make(map[uint32]*subscriber),
	}
}

func (m *msgTopic) Publish(topic string, msg interface{}) {
	m.mx.RLock()
	for _, sub := range m.subs {
		sub.queue <- &channelMsg{
			topic: topic,
			msg:   msg,
		}
	}
	m.mx.RUnlock()
}

func (m *msgTopic) Subscribe(queueSize int, fn Handler) (subscribeId uint32) {
	sub := &subscriber{
		handler: fn,
		queue:   make(chan *channelMsg, queueSize),
	}

	go func() {
		for msg := range sub.queue {
			sub.handler(msg.topic, msg.msg)
		}
	}()

	subId := atomic.LoadUint32(&m.subId)

	m.mx.Lock()
	m.subs[subId] = sub
	m.mx.Unlock()
	return subId
}

func (m *msgTopic) Unsubscribe(subscribeId uint32) {
	m.mx.Lock()
	sub, ok := m.subs[subscribeId]
	if ok {
		close(sub.queue)
		delete(m.subs, subscribeId)
	}
	m.mx.Unlock()
}

func (m *msgTopic) Close() {
	m.mx.Lock()
	for _, sub := range m.subs {
		close(sub.queue)
	}

	// 如果不清除, 在调用 Publish 会导致panic
	m.subs = make(map[uint32]*subscriber)

	m.mx.Unlock()
}
