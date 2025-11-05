# 消息总线

---

# 获得

`go get -u github.com/zlyuancn/zmsgbus`

# 示例

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/zlyuancn/zmsgbus"
)

func main() {
	// 订阅前发消息是无意义的
	zmsgbus.Publish(context.Background(), "topic1", "msg")

	// 订阅 topic1
	zmsgbus.Subscribe("topic1", 0, func(ctx context.Context, msg Message) {
		fmt.Println("Subscribe.topic1", msg.Topic(), msg)
	})

	// 订阅 topic2
	zmsgbus.Subscribe("topic2", 0, func(ctx context.Context, msg Message) {
		fmt.Println("Subscribe.topic2", msg.Topic(), msg)
	})

	// 订阅全局
	zmsgbus.SubscribeGlobal(0, func(ctx context.Context, msg Message) {
		fmt.Println("SubscribeGlobal", msg.Topic(), msg)
	})

	zmsgbus.Publish(context.Background(), "topic1", "msg")
	zmsgbus.Publish(context.Background(), "topic2", "msg")

	// 等待消息消费
	time.Sleep(time.Second)

	// 程序关闭前注意要关闭消息中心
	zmsgbus.Close()
}

```
