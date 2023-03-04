package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/realjf/spinlock"
)

func main() {
	d := testLock(4, 100, spinlock.NewSpinLock())
	fmt.Printf("%4.0fms\n", d.Seconds()*1000)
}

func testLock(threads, n int, l sync.Locker) time.Duration {
	var wg sync.WaitGroup
	wg.Add(threads)

	var count1 int
	var count2 int

	start := time.Now()
	for i := 0; i < threads; i++ {
		go func() {
			for i := 0; i < n; i++ {
				l.Lock()
				count1++
				count2 += 2
				l.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	dur := time.Since(start)
	if count1 != threads*n {
		panic("mismatch")
	}
	if count2 != threads*n*2 {
		panic("mismatch")
	}
	return dur
}
