package mergewrite

import (
	"laneIM/src/config"
	"sync"
	"time"
)

const ()

type MergeWriter struct {
	mu               sync.Mutex
	writeFunc        map[string]*Function
	writeFuncWithArg map[string]*FunctionWithArg
	MaxTime          time.Duration //ms
	MaxCount         int           //count
}

type Function struct {
	f      func() (any, error)
	timer  *time.Timer
	count  int
	signCh chan struct{}
	reply  any
	err    error
	doing  bool
}

type FunctionWithArg struct {
	f      func(any) (any, error)
	timer  *time.Timer
	count  int
	signCh chan struct{}
	arg    any
	reply  any
	err    error
	doing  bool
}

func NewMergeWriter(conf config.BatchWriter) *MergeWriter {
	return &MergeWriter{
		writeFunc:        make(map[string]*Function),
		writeFuncWithArg: make(map[string]*FunctionWithArg),
		MaxTime:          time.Duration(conf.MaxTime) * time.Millisecond,
		MaxCount:         conf.MaxCount,
	}
}

func (m *MergeWriter) Do(key string, f func() (any, error)) (any, error) {
	m.mu.Lock()
	if _, exist := m.writeFunc[key]; !exist {

		// start the timer
		m.writeFunc[key] = &Function{
			timer:  time.NewTimer(m.MaxTime),
			count:  0,
			signCh: make(chan struct{}),
		}

	}
	// update the last f
	function := m.writeFunc[key]
	function.f = f
	function.count++
	if function.count != m.MaxCount {
		m.mu.Unlock()
	} else {
		// reach MaxCount , then do tht f
		delete(m.writeFunc, key)
		function.doing = true
		m.mu.Unlock()

		// do func
		function.reply, function.err = function.f()
		for range function.count - 1 {
			// wake up others(in counter)
			function.signCh <- struct{}{}
		}
		return function.reply, function.err
	}
wait:
	// block untill the last f return
	select {
	case <-function.timer.C:
		m.mu.Lock()
		if function.doing {
			m.mu.Unlock()
			goto wait
		}
		delete(m.writeFunc, key)
		function.doing = true
		m.mu.Unlock()
		function.reply, function.err = function.f()
		for range function.count - 1 {
			// wake up others(in timer)
			function.signCh <- struct{}{}
		}
		// func complete, return
	case <-function.signCh:
		// func complete, return
	}
	return function.reply, function.err
}

func (m *MergeWriter) DoWithArg(key string, arg any, f func(any) (any, error)) (any, error) {
	m.mu.Lock()
	if _, exist := m.writeFuncWithArg[key]; !exist {

		// start the timer
		m.writeFuncWithArg[key] = &FunctionWithArg{
			arg:    arg,
			timer:  time.NewTimer(m.MaxTime),
			count:  0,
			signCh: make(chan struct{}),
		}

	}
	// update the last f
	function := m.writeFuncWithArg[key]
	function.f = f
	function.arg = arg
	function.count++
	if function.count != m.MaxCount {
		m.mu.Unlock()
	} else {
		// reach MaxCount , then do tht f
		delete(m.writeFuncWithArg, key)
		function.doing = true
		m.mu.Unlock()

		// do func
		function.reply, function.err = function.f(function.arg)
		for range function.count - 1 {
			// wake up others(in counter)
			function.signCh <- struct{}{}
		}
		return function.reply, function.err
	}
wait:
	// block untill the last f return
	select {
	case <-function.timer.C:
		m.mu.Lock()
		if function.doing {
			m.mu.Unlock()
			goto wait
		}
		delete(m.writeFuncWithArg, key)
		function.doing = true
		m.mu.Unlock()
		function.reply, function.err = function.f(function.arg)
		for range function.count - 1 {
			// wake up others(in timer)
			function.signCh <- struct{}{}
		}
		// func complete, return
	case <-function.signCh:
		// func complete, return
	}
	return function.reply, function.err
}
