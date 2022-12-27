package gopoolx_test

import (
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/unsafe-risk/utilx/poolx/gopoolx"
)

func TestTimedPoolLeak(t *testing.T) {
	var v uint64
	pool, err := gopoolx.NewTimedPool(
		math.MaxInt64,
		func(v *uint64) {
			atomic.AddUint64(v, 1)
		},
		time.Millisecond*500,
		time.Millisecond*100,
		0,
	)
	if err != nil {
		panic(err)
	}
	defer pool.Stop()

	const N = 100000
	for i := 0; i < N; i++ {
		ok := pool.Run(&v)
		if !ok {
			panic("pool is full")
		}
	}
	time.Sleep(time.Second * 2)

	if atomic.LoadUint64(&v) != N {
		t.Errorf("v = %d, want %d", v, N)
	}

	if w := pool.Workers(); w != 0 {
		t.Errorf("workers = %d, want 0, goroutine leak!!!", w)
	}
}

func TestTimedPoolStop(t *testing.T) {
	var v uint64
	pool, err := gopoolx.NewTimedPool(
		math.MaxInt64,
		func(v *uint64) {
			atomic.AddUint64(v, 1)
		},
		time.Hour,
		time.Millisecond*100,
		0,
	)
	if err != nil {
		panic(err)
	}

	const N = 100000
	for i := 0; i < N; i++ {
		ok := pool.Run(&v)
		if !ok {
			panic("pool is full")
		}
	}
	pool.Stop() // stop the pool
	time.Sleep(time.Second * 2)

	if atomic.LoadUint64(&v) != N {
		t.Errorf("v = %d, want %d", v, N)
	}

	if w := pool.Workers(); w != 0 {
		t.Errorf("workers = %d, want 0, goroutine leak!!!", w)
	}
}

func TestTimedPoolPreheat(t *testing.T) {
	pool, err := gopoolx.NewTimedPool(
		math.MaxInt64,
		func(v *uint64) {
			atomic.AddUint64(v, 1)
		},
		time.Hour,
		time.Millisecond*100,
		10,
	)
	if err != nil {
		panic(err)
	}

	if w := pool.Workers(); w != 10 {
		t.Errorf("workers = %d, want 10", w)
	}

	pool.Stop()
	time.Sleep(time.Second * 1)

	if w := pool.Workers(); w != 0 {
		t.Errorf("workers = %d, want 0, goroutine leak!!!", w)
	}
}

func testIOJob30ms(wg *sync.WaitGroup) {
	time.Sleep(time.Millisecond * 30)
	wg.Done()
}

func BenchmarkTimedPool(b *testing.B) {
	pool, err := gopoolx.NewTimedPool(
		math.MaxInt64,
		testIOJob30ms,
		time.Second*30,
		time.Minute,
		0,
	)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		pool.Run(&wg)
	}
	workers := pool.Workers()
	b.ReportMetric(float64(workers), "workers")
	wg.Wait()
	b.StopTimer()

	pool.Stop()
}

func BenchmarkTimedPoolPar(b *testing.B) {
	pool, err := gopoolx.NewTimedPool(
		math.MaxInt64,
		testIOJob30ms,
		time.Second*30,
		time.Minute,
		0,
	)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			wg.Add(1)
			pool.Run(&wg)
		}
	})
	workers := pool.Workers()
	b.ReportMetric(float64(workers), "workers")
	wg.Wait()
	b.StopTimer()

	pool.Stop()
}

func TestTimedPoolWorkerCount(t *testing.T) {
	pool, err := gopoolx.NewTimedPool(
		math.MaxInt64,
		testIOJob30ms,
		time.Second*30,
		time.Minute,
		0,
	)
	var wg sync.WaitGroup

	if err != nil {
		panic(err)
	}

	if w := pool.Workers(); w != 0 {
		t.Errorf("workers = %d, want 0", w)
	}

	wg.Add(1)
	pool.Run(&wg)
	wg.Wait()

	if w := pool.Workers(); w != 1 {
		t.Errorf("workers = %d, want 1", w)
	}

	wg.Add(2)
	pool.Run(&wg)
	pool.Run(&wg)
	wg.Wait()

	if w := pool.Workers(); w != 2 {
		t.Errorf("workers = %d, want 2", w)
	}

	// Stop the pool
	pool.Stop()
	time.Sleep(time.Second * 1)

	if w := pool.Workers(); w != 0 {
		t.Errorf("workers = %d, want 0, goroutine leak!!!", w)
	}
}
