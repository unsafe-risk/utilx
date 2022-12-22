package sleepx

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/unsafe-risk/utilx/timex"
)

type Sleeper struct {
	// wakeup signal
	ctx    context.Context
	cancel context.CancelFunc

	// slept
	slept atomic.Bool
}

func New() *Sleeper {
	ctx, cancel := context.WithCancel(context.Background())
	return &Sleeper{
		ctx:    ctx,
		cancel: cancel,
		slept:  atomic.Bool{},
	}
}

func (s *Sleeper) SleepFor(d time.Duration) bool {
	if !s.slept.CompareAndSwap(false, true) {
		return false
	}
	defer s.slept.Store(false)
	done := s.ctx.Done()
	for {
		select {
		case <-time.After(time.Duration(d) * time.Second):
			return true
		case <-done:
			return true
		}
	}
}

func (s *Sleeper) SleepUntil(t time.Time) bool {
	if !s.slept.CompareAndSwap(false, true) {
		return false
	}
	defer s.slept.Store(false)
	done := s.ctx.Done()
	for {
		select {
		case <-time.After(t.Sub(timex.Now())):
			return true
		case <-done:
			return true
		}
	}
}

func (s *Sleeper) SleepForever() bool {
	if !s.slept.CompareAndSwap(false, true) {
		return false
	}
	defer s.slept.Store(false)
	done := s.ctx.Done()
	for {
		select {
		case <-done:
			return true
		}
	}
}

func (s *Sleeper) Wakeup() bool {
	if !s.slept.Load() {
		return false
	}
	s.cancel()
	s.ctx, s.cancel = context.WithCancel(context.Background())
	return true
}
