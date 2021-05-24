package zmsgbus

// 处理函数
type ProcessFn func(topic string, msg interface{})

// 订阅者
type subscriber struct {
	fn    ProcessFn
	queue chan interface{}
}
