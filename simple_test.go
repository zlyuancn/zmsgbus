package zmsgbus

import (
	"fmt"
	"sync"
	"testing"
)

func TestSimple(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(4)

	Subscribe("topic1", 1, func(topic string, msg interface{}) {
		fmt.Println("Subscribe.topic1", topic, msg)
		wg.Done()
	})

	Subscribe("topic2", 1, func(topic string, msg interface{}) {
		fmt.Println("Subscribe.topic2", topic, msg)
		wg.Done()
	})

	SubscribeGlobal(1, func(topic string, msg interface{}) {
		fmt.Println("SubscribeGlobal", topic, msg)
		wg.Done()
	})

	Publish("topic1", "msg")
	Publish("topic2", "msg")

	wg.Wait()
}
