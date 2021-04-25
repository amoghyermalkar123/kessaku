package kessaku

import (
	"sync"
	"testing"
	"time"
)

const (
	RunTimes   = 1000000
	BenchParam = 10
	// BenchAntsSize      = 200000
	// DefaultExpiredTime = 10 * time.Second
)

func demoFunc() {
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
}

func BenchmarkKessakuPool(b *testing.B) {
	var wg sync.WaitGroup
	p, _ := NewPool(WithContext(false), WithPoolSize(5))

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			_ = p.Submit(func() {
				demoFunc()
				wg.Done()
			})
		}
		wg.Wait()
	}
	b.StopTimer()
}
