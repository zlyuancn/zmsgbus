package zmsgbus

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestSimple(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(4) // global两次

	Publish(context.Background(), "topic1", "msg")

	Subscribe("topic1", 0, func(ctx context.Context, msg Message) {
		fmt.Println("Subscribe.topic1", msg.Topic(), msg)
		wg.Done()
	})

	Subscribe("topic2", 0, func(ctx context.Context, msg Message) {
		fmt.Println("Subscribe.topic2", msg.Topic(), msg)
		wg.Done()
	})

	// 全局订阅
	SubscribeGlobal(0, func(ctx context.Context, msg Message) {
		fmt.Println("SubscribeGlobal", msg.Topic(), msg)
		wg.Done()
	})

	Publish(context.Background(), "topic1", "msg")
	Publish(context.Background(), "topic2", "msg")

	wg.Wait()
}
