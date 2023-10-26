package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

type MyLockstruct struct {
	mu sync.Mutex
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	newContext, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	for i := 0; i < 2; i++ {
		go func(ctx context.Context, i int) {
			ticker := time.NewTicker(500 * time.Millisecond)

			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Stopping worker %v\n", i)
					return
				case tickTime := <-ticker.C:
					var myMute = MyLockstruct{}
					myMute.mu.Lock()
					defer myMute.mu.Unlock()
					fmt.Printf("Worker %v is ticking at time:%v\n", i, tickTime)
				}
			}
		}(newContext, i)
	}

	time.Sleep(3 * time.Second)
}
