package zmsgbus

var defaultMsgBus = NewMsgBus()

// 发布
func Publish(topic string, msg interface{}) {
	defaultMsgBus.Publish(topic, msg)
}

// 订阅, 返回订阅号
func Subscribe(topic string, fn Handler) (subscribeId uint32) {
	return defaultMsgBus.Subscribe(topic, fn)
}

// 全局订阅, 会收到所有消息, 返回订阅号
func SubscribeGlobal(fn Handler) (subscribeId uint32) {
	return defaultMsgBus.SubscribeGlobal(fn)
}

// 取消订阅
func Unsubscribe(topic string, subscribeId uint32) {
	defaultMsgBus.Unsubscribe(topic, subscribeId)
}

// 取消全局订阅
func UnsubscribeGlobal(subscribeId uint32) {
	defaultMsgBus.UnsubscribeGlobal(subscribeId)
}

// 关闭主题, 同时关闭所有订阅该主题的订阅者
func Close(topic string) {
	defaultMsgBus.Close(topic)
}
