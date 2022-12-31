package decouplex

import (
	"runtime"
	"sync/atomic"
)

type Decoupler struct {
	_read  uint64
	_      [15]uint64
	_write uint64
	_      [15]uint64

	size    uint64
	locktab []uint64
	buffer  [][]byte
	drop    func(b []byte) (int, error)
}

func NewDecoupler(size uint64, drop func(b []byte) (int, error)) *Decoupler {
	if drop == nil {
		drop = func(b []byte) (int, error) {
			return len(b), nil
		}
	}

	locktab := make([]uint64, size)
	for i := range locktab {
		locktab[i] = empty
	}

	return &Decoupler{
		size:    size,
		_write:  ^uint64(0),
		locktab: locktab,
		buffer:  make([][]byte, size),
		drop:    drop,
	}
}

const (
	empty = ^uint64(1)
	busy  = ^uint64(0)
)

func (d *Decoupler) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return d.drop(b)
	}

	for {
		ts := atomic.AddUint64(&d._write, 1)
		idx := ts % d.size

		lock := atomic.LoadUint64(&d.locktab[idx])
		if lock == empty {
			if !atomic.CompareAndSwapUint64(&d.locktab[idx], empty, busy) {
				continue
			}
			d.buffer[idx] = append(d.buffer[idx][:0], b...)
			atomic.StoreUint64(&d.locktab[idx], ts)
			return len(b), nil
		}

		if lock != busy {
			r := atomic.LoadUint64(&d._read)
			if ts-r >= d.size {
				if lock < ts {
					return d.drop(b)
				}
			}
		}
	}
}

func (d *Decoupler) Next(view func([]byte)) bool {
	w := atomic.LoadUint64(&d._write)
	r := atomic.LoadUint64(&d._read)
	if w <= r {
		return false
	}
	idx := r % d.size

	var lock uint64
	for {
		lock = atomic.LoadUint64(&d.locktab[idx])
		if lock == empty || lock == busy {
			runtime.Gosched() // Wait for the writer
			continue
		}
		break
	}

	if view != nil {
		view(d.buffer[idx])
	}
	atomic.StoreUint64(&d.locktab[idx], empty)
	atomic.AddUint64(&d._read, 1)

	return true
}
