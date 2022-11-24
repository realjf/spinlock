# spinlock

reentrant spin lock(可重入自旋锁)

## Usage
```go
lock := NewSpinLock()
lock.Lock()
defer lock.Unlock()

// or
lock.TryLock()

// or
lock.SetTryTimes(2)
lock.TryLock()
```
