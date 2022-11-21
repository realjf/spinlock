package spinlock

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func TestSpinLock(t *testing.T) {
	// 调整vscode的配置go.testTimeout时间
	cases := []struct {
		name    string
		threads int
		n       int
		l       sync.Locker
	}{
		{
			name:    "spinlock[1]",
			threads: 1,
			n:       1000000,
			l:       NewSpinLock(),
		},
		{
			name:    "mutex[1]",
			threads: 1,
			n:       1000000,
			l:       &sync.Mutex{},
		},
		{
			name:    "spinlock[4]",
			threads: 4,
			n:       1000000,
			l:       NewSpinLock(),
		},
		{
			name:    "mutex[4]",
			threads: 4,
			n:       1000000,
			l:       &sync.Mutex{},
		},
		{
			name:    "spinlock[8]",
			threads: 8,
			n:       1000000,
			l:       NewSpinLock(),
		},
		{
			name:    "mutex[8]",
			threads: 8,
			n:       1000000,
			l:       &sync.Mutex{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ti := testLock(tc.threads, tc.n, tc.l)
			assert.NotNil(t, ti)
			t.Logf("%s %4.0fms\n", tc.name, ti.Seconds()*1000)
		})
	}
}
