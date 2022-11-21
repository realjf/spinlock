package spinlock

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type SpinLock struct {
	owner int64
	count int
}

func (sl *SpinLock) get() *int64 {
	return &sl.owner
}

func (sl *SpinLock) Lock() {
	me := GetGoroutineId()
	if sl.owner == me { // 如果当前线程已经获取到了锁，线程数增加一，然后返回
		sl.count++
		return
	}
	// 如果没获取到锁，则通过CAS自旋
	for !atomic.CompareAndSwapInt64(sl.get(), 0, me) {
		runtime.Gosched()
	}
}
func (sl *SpinLock) Unlock() {
	if sl.owner != GetGoroutineId() {
		panic("illegalMonitorStateError")
	}
	if sl.count > 0 { // 如果大于0，表示当前线程多次获取了该锁，释放锁通过count减一来模拟
		sl.count--
	} else { // 如果count==0，可以将锁释放，这样就能保证获取锁的次数与释放锁的次数是一致的了。
		atomic.StoreInt64(sl.get(), 0)
	}
}

func GetGoroutineId() int64 {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic recover:panic info:%v\n", err)
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.ParseInt(idField, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

func NewSpinLock() sync.Locker {
	var lock SpinLock
	return &lock
}
