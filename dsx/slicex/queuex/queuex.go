package queuex

import (
	"github.com/unsafe-risk/utilx/dsx"
)

var _ = (dsx.Queue[int])(&Ring[int]{})

type Ring[T any] struct {
	data []T
	cap  int
	head int
	tail int
}

func New[T any](cap int) *Ring[T] {
	if cap == 0 {
		return nil
	}
	return &Ring[T]{
		data: make([]T, cap),
		cap:  cap + 1,
		head: 0,
		tail: 0,
	}
}

func (r *Ring[T]) Enqueue(data T) bool {
	if r.IsFull() == true {
		return false
	}
	r.data[r.tail] = data
	r.tail = (r.tail + 1) % r.cap
	return true
}

func (r *Ring[T]) Dequeue() (data T, ok bool) {
	if r.IsEmpty() {
		return
	}
	data = r.data[r.head]
	r.head = (r.head + 1) % r.cap
	return data, true
}

func (r *Ring[T]) Head() T {
	return r.data[r.head]
}

func (r *Ring[T]) Cap() int {
	return r.cap - 1
}

func (r *Ring[T]) IsEmpty() bool {
	return r.head == r.tail
}

func (r *Ring[T]) IsFull() bool {
	return r.head == (r.tail+1)%r.cap
}

func (r *Ring[T]) Peek() (data T, ok bool) {
	if r.IsEmpty() {
		return
	}
	return r.data[r.head], true
}
