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
	owner     int64
	count     int
	try_times int // 0 indicates that you will try to lock until you obtain.
}

func (sl *SpinLock) get() *int64 {
	return &sl.owner
}

// 0 indicates that you will try to lock until you obtain.
func (sl *SpinLock) SetTryTimes(n int) bool {
	sl.try_times = n
	return true
}

// true - locked
// false - unlocked
func (sl *SpinLock) TryLock() bool {
	me := GetGoroutineId()
	if sl.owner == me {
		// If the current thread has obtained the lock, increase the number of threads by one, and then return
		sl.count++
		return true
	}

	return atomic.CompareAndSwapInt64(sl.get(), 0, me)
}

func (sl *SpinLock) Lock() {
	// If the lock is not obtained, spin it through CAS
	if sl.try_times > 0 {
		for loop := 0; loop < sl.try_times; loop++ {
			if !sl.TryLock() {
				runtime.Gosched()
			} else {
				break
			}
		}
	} else {
		for !sl.TryLock() {
			runtime.Gosched()
		}
	}
}
func (sl *SpinLock) Unlock() {
	if sl.owner != GetGoroutineId() {
		panic("illegalMonitorStateError")
	}
	if sl.count > 0 {
		// If it is greater than 0, it means that the current thread has acquired the lock multiple times,\
		// and releasing the lock is simulated by subtracting count
		sl.count--
	} else {
		// If count==0, the lock can be released, so that the number of times \
		// to obtain the lock is consistent with the number of times to release the lock.
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
