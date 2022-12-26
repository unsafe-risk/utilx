package poolx

import "sync"

// Pool is wrapper around sync.Pool.
//
// It provides a type-safe interface for sync.Pool.
type Pool[T any] struct {
	pool sync.Pool
}

// Init initializes the pool.
// It MUST be called before using the pool.
//
// Note: This function is not thread-safe.
func (p *Pool[T]) Init(New func() *T) {
	p.pool.New = func() interface{} {
		return New()
	}
}

// Get returns an object from the pool.
//
// Note: This function is thread-safe.
func (p *Pool[T]) Get() *T {
	return p.pool.Get().(*T)
}

// Put returns an object to the pool.
// The object MUST not be used after this call.
// The object MUST reset itself before returning to the pool.
//
// Note: This function is thread-safe.
func (p *Pool[T]) Put(x *T) {
	p.pool.Put(x)
}
