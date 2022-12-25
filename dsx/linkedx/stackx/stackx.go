package stackx

import (
	"sync"

	"github.com/unsafe-risk/utilx/dsx"
)

type node[T any] struct {
	data T
	next *node[T]
}

var _ = (dsx.Stack[int])(&Stack[int]{})

type Stack[T any] struct {
	head   *node[T]
	length int
	pool   *sync.Pool
}

func New[T any]() *Stack[T] {
	return &Stack[T]{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(node[T])
			},
		},
	}
}

func (s *Stack[T]) Push(data T) {
	s.length++
	n := s.pool.Get().(*node[T])
	n.data = data
	n.next = s.head
	s.head = n
}

func (s *Stack[T]) Pop() (data T, ok bool) {
	if s.head == nil {
		return
	}
	s.length--
	s.pool.Put(s.head)
	s.head = s.head.next
	return s.head.data, true
}

func (s *Stack[T]) Peek() (data T, ok bool) {
	if s.head == nil {
		return
	}
	return s.head.data, true
}

func (s *Stack[T]) Len() int {
	return s.length
}
