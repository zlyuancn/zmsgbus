package zmsgbus

import (
	"context"
	"sync/atomic"
)

type Subscriber interface {
	// 获取订阅者id
	GetSubId() uint32
	// 启动
	Start()
	// 关闭订阅者, 程序结束前注意要调用这个方法
	Close()

	Handler(msg Message)
}

// 全局自增订阅者id
var autoIncrSubscriberId uint32

// 生成下一个订阅者id
func nextSubscriberId() uint32 {
	return atomic.AddUint32(&autoIncrSubscriberId, 1)
}

// 处理函数
type Handler func(ctx context.Context, msg Message)

// 订阅者
type subscriber struct {
	subId       uint32
	handler     Handler
	queue       chan Message
	threadCount int

	onceStart int32
	onceClose int32
}

func newSubscriber(msgQueueSize int, threadCount int, handler Handler) Subscriber {
	sub := &subscriber{
		subId:   nextSubscriberId(),
		handler: handler,
		queue:   make(chan Message, msgQueueSize),
	}

	if threadCount < 1 {
		threadCount = 1
	}
	sub.threadCount = threadCount
	return sub
}

func (s *subscriber) GetSubId() uint32 {
	return s.subId
}
func (s *subscriber) Start() {
	if atomic.AddInt32(&s.onceStart, 1) != 1 {
		return
	}

	for i := 0; i < s.threadCount; i++ {
		go s.start()
	}
}
func (s *subscriber) start() {
	for msg := range s.queue {
		s.handler(msg.Ctx(), msg)
	}
}

func (s *subscriber) Close() {
	if atomic.AddInt32(&s.onceClose, 1) == 1 {
		close(s.queue)
	}
}

func (s *subscriber) Handler(msg Message) {
	s.queue <- msg
}
