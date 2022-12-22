package ringx

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

func (t *Ring[T]) Enqueue(data T) error {
	if t.IsFull() == true {
		return fmt.Errorf("Queue is full | cap - %d", t.cap)
	}
	t.data[t.tail] = data
	t.tail = (t.tail + 1) % t.cap
	return nil
}

func (t *Ring[T]) Dequeue() (data T, ok bool) {
	if t.IsEmpty() {
		return
	}
	data = t.data[t.head]
	t.head = (t.head + 1) % t.cap
	return data, true
}

func (t *Ring[T]) Head() T {
	return t.data[t.head]
}

func (t *Ring[T]) Cap() int64 {
	return t.cap
}

func (t *Ring[T]) IsEmpty() bool {
	return t.head == t.tail
}

func (t *Ring[T]) IsFull() bool {
	return t.head == (t.tail+1)%t.cap
}
