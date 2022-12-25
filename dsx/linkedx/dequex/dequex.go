package dequex

import "sync"

const nodeBufferSize = 64

type node[T any] struct {
	data  [nodeBufferSize]T
	start int
	end   int
	next  *node[T]
	prev  *node[T]
}

type Deque[T any] struct {
	pool sync.Pool
	head *node[T]
	tail *node[T]
}

func New[T any]() *Deque[T] {
	return &Deque[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return new(node[T])
			},
		},
	}
}

func (t *Deque[T]) PushFront(data T) {
	if t.head == nil {
		n := t.pool.Get().(*node[T])
		n.data[0] = data
		n.start = 0
		n.end = 1
		t.head = n
		t.tail = n
		return
	}
	if t.head.start == 0 {
		n := t.pool.Get().(*node[T])
		n.data[nodeBufferSize-1] = data
		n.start = nodeBufferSize - 1
		n.end = nodeBufferSize
		n.next = t.head
		t.head.prev = n
		t.head = n
		return
	}
	t.head.start--
	t.head.data[t.head.start] = data
}

func (t *Deque[T]) PushBack(data T) {
	if t.tail == nil {
		n := t.pool.Get().(*node[T])
		n.data[0] = data
		n.start = 0
		n.end = 1
		t.head = n
		t.tail = n
		return
	}
	if t.tail.end == nodeBufferSize {
		n := t.pool.Get().(*node[T])
		n.data[0] = data
		n.start = 0
		n.end = 1
		n.prev = t.tail
		t.tail.next = n
		t.tail = n
		return
	}
	t.tail.data[t.tail.end] = data
	t.tail.end++
}

func (t *Deque[T]) PopFront() (data T, ok bool) {
	if t.head == nil {
		return
	}
	data = t.head.data[t.head.start]
	t.head.start++
	if t.head.start == t.head.end {
		t.pool.Put(t.head)
		t.head = t.head.next
		if t.head != nil {
			t.head.prev = nil
		}
	}
	return data, true
}

func (t *Deque[T]) PopBack() (data T, ok bool) {
	if t.tail == nil {
		return
	}
	t.tail.end--
	data = t.tail.data[t.tail.end]
	if t.tail.start == t.tail.end {
		t.pool.Put(t.tail)
		t.tail = t.tail.prev
		if t.tail != nil {
			t.tail.next = nil
		}
	}
	return data, true
}

func (t *Deque[T]) PeekFront() (data T, ok bool) {
	if t.head == nil {
		return
	}
	return t.head.data[t.head.start], true
}

func (t *Deque[T]) PeekBack() (data T, ok bool) {
	if t.tail == nil {
		return
	}
	return t.tail.data[t.tail.end-1], true
}
