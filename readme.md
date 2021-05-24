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
	wg.Add(2)

	zmsgbus.Subscribe("topic", func(topic string, msg interface{}) {
		fmt.Println(msg)
		wg.Done()
	})

	zmsgbus.Subscribe("topic", func(topic string, msg interface{}) {
		fmt.Println(msg)
		wg.Done()
	})

	zmsgbus.Publish("topic", "msg")

	wg.Wait()
}
```
