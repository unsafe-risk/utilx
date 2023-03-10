package queuex

import (
	"sync"

	"github.com/unsafe-risk/utilx/dsx"
)

type node[T any] struct {
	data T
	next *node[T]
}

var _ = (dsx.Queue[int])(&Queue[int]{})

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

func (q *Queue[T]) Enqueue(data T) bool {
	n := q.pool.Get().(*node[T])
	n.data = data
	if q.head == nil {
		q.head = n
		q.tail = n
		return true
	}
	q.tail.next = n
	q.tail = n
	return true
}

func (q *Queue[T]) Dequeue() (data T, ok bool) {
	if q.head == nil {
		return
	}
	q.pool.Put(q.head)
	q.head = q.head.next
	return q.head.data, true
}

func (q *Queue[T]) Peek() (data T, ok bool) {
	if q.head == nil {
		return
	}
	return q.head.data, true
}

func (q *Queue[T]) IsEmpty() bool {
	return q.head == nil
}
