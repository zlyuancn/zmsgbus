package zmsgbus

import (
	"fmt"
	"sync"
	"testing"
)

func TestTopic(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(10)

	topic1 := newMsgTopic()
	defer topic1.Close()

	topic2 := newMsgTopic()
	defer topic2.Close()

	subscribe1 := topic1.Subscribe(10, func(topic string, msg interface{}) {
		fmt.Println("subscribe.topic1", msg)
		wg.Done()
	})

	subscribe2 := topic2.Subscribe(10, func(topic string, msg interface{}) {
		fmt.Println("subscribe.topic2", msg)
		wg.Done()
	})

	for i := 0; i < 5; i++ {
		topic1.Publish("topic1", i)
		topic2.Publish("topic2", i)
	}

	topic1.Unsubscribe(subscribe1)
	topic2.Unsubscribe(subscribe2)

	wg.Wait()
}
