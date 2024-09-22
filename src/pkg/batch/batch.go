package batch

import (
	"sync"
	"time"
)

type BatchArgElem interface {
}
type BatchArgs[T BatchArgElem] struct {
	args  []*T
	index int
	mu    sync.RWMutex
	size  int
	f     func([]*T)
	Time  time.Duration
	stop  bool
}

func NewBatchArgs[T BatchArgElem](size int, time time.Duration, f func([]*T)) *BatchArgs[T] {
	rt := &BatchArgs[T]{
		args:  make([]*T, size),
		index: 0,
		size:  size,
		f:     f,
		Time:  time,
	}
	return rt
}

func (b *BatchArgs[T]) Add(elem *T) {
Retry:
	b.mu.Lock()
	if b.index == b.size {
		// TODO 重置定时器
		b.mu.Unlock()
		b.Do()
		goto Retry
	}
	b.args[b.index] = elem
	b.index++
	b.mu.Unlock()
}

func (b *BatchArgs[T]) Take() []*T {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.index == 0 {
		return nil
	}
	rt := b.args[:b.index]
	b.args = make([]*T, b.size)
	b.index = 0
	return rt
}

func (b *BatchArgs[T]) Start() {
	go func() {
		for {
			if b.stop {
				break
			}
			go b.Do()
			time.Sleep(b.Time)
		}
	}()
}

func (b *BatchArgs[T]) Do() {
	if rt := b.Take(); rt != nil {
		b.f(rt)
	}
}

func (b *BatchArgs[T]) Stop() {
	b.stop = true
}
