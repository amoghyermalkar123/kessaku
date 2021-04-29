package kessaku

import (
	"sync"
	"testing"
	"time"
)

const (
	runTimes   = 1000000
	benchParam = 10
	// BenchAntsSize      = 200000
	// DefaultExpiredTime = 10 * time.Second
)

func demoFunc() {
	time.Sleep(time.Duration(benchParam) * time.Millisecond)
}

func BenchmarkKessakuPool(b *testing.B) {
	var wg sync.WaitGroup
	p, _ := NewPool()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(runTimes)
		for j := 0; j < runTimes; j++ {
			_ = p.Submit(func() {
				demoFunc()
				wg.Done()
			})
		}
		wg.Wait()
	}
	b.StopTimer()
}
