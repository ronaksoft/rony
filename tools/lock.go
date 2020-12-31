package tools

import (
	"runtime"
	"sync"
	"sync/atomic"
)

/*
   Creation Time: 2021 - Jan - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// SpinLock is a spinlock implementation.
//
// A SpinLock must not be copied after first use.
// This SpinLock intended to be used to synchronize exceptionally short lived operations.
type SpinLock struct {
	_    sync.Mutex // for copy protection compiler warning
	lock uintptr
}

// Lock locks l.
// If the lock is already in use, the calling goroutine
// blocks until the locker is available.
func (l *SpinLock) Lock() {
	for !atomic.CompareAndSwapUintptr(&l.lock, 0, 1) {
		runtime.Gosched()
	}
}

// Unlock unlocks l.
func (l *SpinLock) Unlock() {
	atomic.StoreUintptr(&l.lock, 0)
}
