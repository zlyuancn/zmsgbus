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

type Topic interface {
	// 发布
	Publish(ctx context.Context, topic string, msg interface{})
	// 订阅, 返回订阅号
	Subscribe(msgQueueSize int, threadCount int, handler Handler) (subscribeId uint32)
	// 取消订阅
	Unsubscribe(subscribeId uint32)
	// 关闭主题, 程序结束前注意要调用这个方法
	Close()
}

// 主题
type msgTopic struct {
	subs map[uint32]Subscriber
	mx   sync.RWMutex // 用于锁 subs
}

func newMsgTopic() Topic {
	return &msgTopic{
		subs: make(map[uint32]Subscriber),
	}
}

func (m *msgTopic) Publish(ctx context.Context, topic string, msg interface{}) {
	m.mx.RLock()
	for _, sub := range m.subs {
		sub.Handler(newMessage(ctx, topic, msg))
	}
	m.mx.RUnlock()
}

func (m *msgTopic) Subscribe(msgQueueSize int, threadCount int, handler Handler) (subscribeId uint32) {
	sub := newSubscriber(msgQueueSize, threadCount, handler)
	subId := sub.GetSubId()
	sub.Start()

	m.mx.Lock()
	m.subs[subId] = sub
	m.mx.Unlock()
	return subId
}

func (m *msgTopic) Unsubscribe(subscribeId uint32) {
	m.mx.Lock()
	sub, ok := m.subs[subscribeId]
	if ok {
		delete(m.subs, subscribeId)
	}
	m.mx.Unlock()

	// 后置关闭
	if ok {
		sub.Close()
	}
}

func (m *msgTopic) Close() {
	m.mx.Lock()
	clearSubs := make([]Subscriber, 0, len(m.subs))
	for _, sub := range m.subs {
		clearSubs = append(clearSubs, sub)
	}
	m.subs = make(map[uint32]Subscriber)
	m.mx.Unlock()

	// 后置关闭
	for _, sub := range clearSubs {
		sub.Close()
	}
}
