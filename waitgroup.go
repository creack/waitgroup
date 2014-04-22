package waitgroup

import (
	"sync"
	"sync/atomic"
)

type WaitGroup struct {
	m       sync.Mutex
	counter int32
	waiters int32
	ch      chan struct{}
}

func (wg *WaitGroup) Add(delta int) {
	wg.m.Lock()
	if wg.ch == nil {
		wg.ch = make(chan struct{})
	}
	wg.m.Unlock()
	v := atomic.AddInt32(&wg.counter, int32(delta))
	// Panic if negative
	if v < 0 {
		panic("sync: negative WaitGroup counter")
	} else if v == 0 {
		for waiters := atomic.LoadInt32(&wg.waiters); waiters > 0; waiters-- {
			wg.ch <- struct{}{}
		}
	}
}

func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

func (wg *WaitGroup) Wait() {
	if atomic.LoadInt32(&wg.counter) == 0 {
		return
	}
	wg.m.Lock()
	if wg.ch == nil {
		wg.ch = make(chan struct{})
	}
	wg.m.Unlock()
	atomic.AddInt32(&wg.waiters, 1)
	<-wg.ch
	atomic.AddInt32(&wg.waiters, -1)
}

func main2() {
	_ = sync.WaitGroup{}
}
