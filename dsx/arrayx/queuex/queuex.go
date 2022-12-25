package queuex

import "fmt"

type Ring[T any] struct {
	data []T
	cap  int64
	head int64
	tail int64
}

func New[T any](cap int64) *Ring[T] {
	if cap == 0 {
		return nil
	}
	return &Ring[T]{
		data: make([]T, cap),
		cap:  cap,
		head: 0,
		tail: 0,
	}
}

func (r *Ring[T]) Enqueue(data T) error {
	if r.IsFull() == true {
		return fmt.Errorf("Queue is full | cap - %d", r.cap)
	}
	r.data[r.tail] = data
	r.tail = (r.tail + 1) % r.cap
	return nil
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

func (r *Ring[T]) Cap() int64 {
	return r.cap
}

func (r *Ring[T]) IsEmpty() bool {
	return r.head == r.tail
}

func (r *Ring[T]) IsFull() bool {
	return r.head == (r.tail+1)%r.cap
}
