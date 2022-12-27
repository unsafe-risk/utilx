// TimedPool is a pool of goroutines with a idle timeout.

package gopoolx

import (
	"errors"
	"sync"
	"time"

	"github.com/unsafe-risk/utilx/syncx/syncpoolx"
)

type timedworker[T any] struct {
	ch chan T
	tt int64
}

type TimedPool[T any] struct {
	maxWorkers  int64
	handler     func(T)
	stop        chan struct{}
	idleTimeout time.Duration
	gcPeriod    time.Duration
	_           [10]uint64

	mu      sync.Mutex
	workers int64
	pool    syncpoolx.Pool[timedworker[T]]
	dobby   []*timedworker[T]
}

func tworker[T any](pool *TimedPool[T], w *timedworker[T]) {
	for v := range w.ch {
		pool.handler(v)
		if pool.free(w) {
			return
		}
	}

	pool.mu.Lock()
	pool.workers--
	pool.mu.Unlock()
}

func (pool *TimedPool[T]) start(preheat int) {
	pool.pool.Init(func() *timedworker[T] {
		return &timedworker[T]{
			ch: make(chan T, 1),
		}
	})
	pool.stop = make(chan struct{})

	if preheat > 0 {
		pool.workers = int64(preheat)
		for i := 0; i < preheat; i++ {
			w := pool.pool.Get()
			w.tt = time.Now().UnixNano()
			pool.dobby = append(pool.dobby, w)
			go tworker(pool, w)
		}
	}

	go func() {
		ticker := time.NewTicker(pool.gcPeriod)
		defer ticker.Stop()

		var wbuffer []*timedworker[T]

		for {
			select {
			case <-ticker.C:
				pool.mu.Lock()
				if pool.maxWorkers == -1 {
					// Stop All Workers
					for _, w := range pool.dobby {
						close(w.ch)
					}
					pool.mu.Unlock()
					return
				}

				// Stop Idle Workers
				now := time.Now().Add(-pool.idleTimeout).UnixNano()
				dc := len(pool.dobby)
				if dc > 0 {
					l, r := 0, dc
					m := (l + r) >> 1
					for l < r {
						if pool.dobby[m].tt < now {
							l = m + 1
						} else {
							r = m
						}
						m = (l + r) >> 1
					}

					wbuffer = append(wbuffer, pool.dobby[:l]...)
					copy(pool.dobby, pool.dobby[l:])
					pool.dobby = pool.dobby[:dc-l]
				}
				pool.mu.Unlock()

				// Stop Old Workers
				for i := range wbuffer {
					close(wbuffer[i].ch)
					wbuffer[i].tt = 0
					wbuffer[i].ch = nil
					pool.pool.Put(wbuffer[i])
					wbuffer[i] = nil
				}
				wbuffer = wbuffer[:0]

			case <-pool.stop:
				pool.mu.Lock()
				pool.maxWorkers = -1
				for i := range pool.dobby {
					close(pool.dobby[i].ch)
					// Drop All Workers
					// pool.dobby[i].tt = 0
					// pool.dobby[i].ch = nil
					// pool.pool.Put(pool.dobby[i])
					// pool.dobby[i] = nil
				}
				pool.mu.Unlock()
				return
			}
		}
	}()
}

func (pool *TimedPool[T]) free(w *timedworker[T]) (stop bool) {
	w.tt = time.Now().UnixNano()
	pool.mu.Lock()
	if pool.maxWorkers == -1 {
		pool.workers--
		pool.mu.Unlock()
		return true
	}
	pool.dobby = append(pool.dobby, w)
	pool.mu.Unlock()
	return false
}

func (pool *TimedPool[T]) Run(v T) (ok bool) {
	var w *timedworker[T]
	pool.mu.Lock()
	// Check If Pool Is Draining
	if pool.maxWorkers == -1 {
		pool.mu.Unlock()
		return false
	}

	// Check For Dobby
	if len(pool.dobby) > 0 {
		w = pool.dobby[len(pool.dobby)-1]
		pool.dobby = pool.dobby[:len(pool.dobby)-1]
	} else {
		// No Dobby, Create Worker
		if pool.workers < pool.maxWorkers {
			pool.workers++
			w = pool.pool.Get()
			w.ch = make(chan T, 1)
			go tworker(pool, w)
		} else {
			pool.mu.Unlock()
			return false
		}
	}
	pool.mu.Unlock()

	w.ch <- v
	return true
}

var (
	ErrInvalidMaxWorkers  = errors.New("invalid max workers, must be greater than 0")
	ErrInvalidHandler     = errors.New("invalid handler, must not be nil")
	ErrInvalidIdleTimeout = errors.New("invalid idle timeout, must be greater than 0s")
	ErrInvalidGCPeriod    = errors.New("invalid gc period, must be greater than 0s")
)

func NewTimedPool[T any](maxWorkers int64, handler func(T), idleTimeout, gcPeriod time.Duration, preheat int) (*TimedPool[T], error) {
	if maxWorkers <= 0 {
		return nil, ErrInvalidMaxWorkers
	}
	if handler == nil {
		return nil, ErrInvalidHandler
	}
	if idleTimeout <= 0 {
		return nil, ErrInvalidIdleTimeout
	}
	if gcPeriod <= 0 {
		return nil, ErrInvalidGCPeriod
	}

	pool := &TimedPool[T]{
		maxWorkers:  maxWorkers,
		handler:     handler,
		idleTimeout: idleTimeout,
		gcPeriod:    gcPeriod,
	}

	pool.start(preheat)
	return pool, nil
}

func (pool *TimedPool[T]) Stop() {
	if pool.stop == nil {
		return
	}
	pool.stop <- struct{}{}
	pool.stop = nil
}

func (pool *TimedPool[T]) Workers() int64 {
	pool.mu.Lock()
	w := pool.workers
	pool.mu.Unlock()
	return w
}
