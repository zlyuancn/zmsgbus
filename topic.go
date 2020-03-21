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

// 处理函数
type ProcessFn func(msg interface{})

// 订阅者
type subscriber struct {
    fn    ProcessFn
    queue chan interface{}
}

// 主题们
type msgTopics map[string]*msgTopic

// 主题
type msgTopic struct {
    subid uint32
    subs  map[uint32]*subscriber
    mx    sync.RWMutex
}

func newMsgTopic() *msgTopic {
    return &msgTopic{
        subs: make(map[uint32]*subscriber),
    }
}

func (m *msgTopic) Publish(msg interface{}) {
    m.mx.RLock()
    for _, sub := range m.subs {
        sub.queue <- msg
    }
    m.mx.RUnlock()
}

func (m *msgTopic) Subscribe(queueSize int, fn ProcessFn) (subscribeId uint32) {
    sub := &subscriber{
        fn:    fn,
        queue: make(chan interface{}, queueSize),
    }

    go func() {
        for msg := range sub.queue {
            sub.fn(msg)
        }
    }()

    m.mx.Lock()

    id := m.subid
    m.subid++

    m.subs[id] = sub

    m.mx.Unlock()
    return id
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

    // 别动它! 这个可以防止极端情况下的崩溃
    m.subs = make(map[uint32]*subscriber)

    m.mx.Unlock()
}
