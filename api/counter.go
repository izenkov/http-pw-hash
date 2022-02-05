package api

import "sync/atomic"

type AtomicInt struct {
	counter int64
}

func NewAtomicInt() *AtomicInt {
	ai := &AtomicInt{
		counter: 0,
	}
	return ai
}

func (ai *AtomicInt) Inc() int64 {
	return atomic.AddInt64(&ai.counter, 1)
}

func (ai *AtomicInt) Get() int64 {
	return atomic.LoadInt64(&ai.counter)
}

type AtomicAverage struct {
	counter int64
	average int64
}

func NewAtomicAverage() *AtomicAverage {
	aa := &AtomicAverage{
		counter: 0,
		average: 0,
	}
	return aa
}

func (aa *AtomicAverage) Add(i64 int64) {
	oc := atomic.LoadInt64(&aa.counter)
	oa := atomic.LoadInt64(&aa.average)
	nc := oc + 1
	na := ((oa * oc) + i64) / nc
	atomic.StoreInt64(&aa.counter, nc)
	atomic.StoreInt64(&aa.average, na)
}

func (aa *AtomicAverage) Get() (cnt int64, ave int64) {
	return atomic.LoadInt64(&aa.counter), atomic.LoadInt64(&aa.average)
}
