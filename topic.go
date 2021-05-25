/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/3/20
   Description :
-------------------------------------------------
*/

package zmsgbus

import (
	"runtime"
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

func (m *msgTopic) Subscribe(queueSize int, threadCount int, handler Handler) (subscribeId uint32) {
	sub := &subscriber{
		handler: handler,
		queue:   make(chan *channelMsg, queueSize),
	}

	if threadCount < 1 {
		threadCount = runtime.NumCPU() >> 1
		if threadCount < 1 {
			threadCount = 1
		}
	}
	for i := 0; i < threadCount; i++ {
		go m.start(sub)
	}

	subId := atomic.LoadUint32(&m.subId)

	m.mx.Lock()
	m.subs[subId] = sub
	m.mx.Unlock()
	return subId
}

func (m *msgTopic) start(sub *subscriber) {
	for msg := range sub.queue {
		sub.handler(msg.topic, msg.msg)
	}
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
