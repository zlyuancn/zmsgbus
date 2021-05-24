package zmsgbus

// 发布
func Publish(topic string, msg interface{}) {
	defaultMsgBus.Publish(topic, msg)
}

// 订阅, 返回订阅号
func Subscribe(topic string, fn ProcessFn) (subscribeId uint32) {
	return defaultMsgBus.Subscribe(topic, fn)
}

// 取消订阅
func Unsubscribe(topic string, subscribeId uint32) {
	defaultMsgBus.Unsubscribe(topic, subscribeId)
}

// 关闭主题, 同时关闭所有订阅该主题的订阅者
func Close(topic string) {
	defaultMsgBus.Close(topic)
}
