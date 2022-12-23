package lazyx

import "sync"

func New[T any](ctor func() T) *Lazy[T] {
	return &Lazy[T]{
		ctor: ctor,
	}
}

type Lazy[T any] struct {
	o     sync.Once
	value T
	ctor  func() T
}

func (l *Lazy[T]) Get() T {
	l.o.Do(func() {
		l.value = l.ctor()
	})

	return l.value
}
