package utils

import (
	"sync"
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	limiter := NewLimiter(1)
	var num int
	var lock sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			if limiter.Get() {
				lock.Lock()
				num++
				lock.Unlock()

				time.Sleep(time.Microsecond)
				limiter.Put()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	t.Log(num)
}
