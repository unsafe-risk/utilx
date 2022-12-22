package queuex

import "sync"

type node[T any] struct {
	data T
	next *node[T]
}

type Queue[T any] struct {
	pool sync.Pool
	head *node[T]
	tail *node[T]
}

func New[T any]() *Queue[T] {
	return &Queue[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return new(node[T])
			},
		},
	}
}

func (t *Queue[T]) Enqueue(data T) {
	n := t.pool.Get().(*node[T])
	n.data = data
	if t.head == nil {
		t.head = n
		t.tail = n
		return
	}
	t.tail.next = n
	t.tail = n
}

func (t *Queue[T]) Dequeue() (data T, ok bool) {
	if t.head == nil {
		ok = false
		return
	}
	data = t.head.data
	ok = true
	t.pool.Put(t.head)
	t.head = t.head.next
	return
}

func (t *Queue[T]) Peek() (data T, ok bool) {
	if t.head == nil {
		ok = false
		return
	}
	data = t.head.data
	ok = true
	return
}

func (t *Queue[T]) IsEmpty() bool {
	return t.head == nil
}
