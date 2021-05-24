# 消息总线

---

# 获得

`go get -u github.com/zlyuancn/zmsgbus`

# 示例

```go
package main

import (
	"fmt"
	"sync"

	"github.com/zlyuancn/zmsgbus"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(4)

	zmsgbus.Subscribe("topic1", func(topic string, msg interface{}) {
		fmt.Println("Subscribe.topic1", topic, msg)
		wg.Done()
	})

	zmsgbus.Subscribe("topic2", func(topic string, msg interface{}) {
		fmt.Println("Subscribe.topic2", topic, msg)
		wg.Done()
	})

	zmsgbus.SubscribeGlobal(func(topic string, msg interface{}) {
		fmt.Println("SubscribeGlobal", topic, msg)
		wg.Done()
	})

	zmsgbus.Publish("topic1", "msg")
	zmsgbus.Publish("topic2", "msg")

	wg.Wait()
}
```
