package listx

import (
	"sync"

	"github.com/unsafe-risk/utilx/dsx"
)

type node[T any] struct {
	data T
	next *node[T]
}

var _ = (dsx.List[int])(&List[int]{})

type List[T any] struct {
	head   *node[T]
	tail   *node[T]
	length int
	pool   *sync.Pool
}

func New[T any]() *List[T] {
	return &List[T]{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(node[T])
			},
		},
	}
}

func (l *List[T]) Append(data T) {
	l.length++
	n := l.pool.Get().(*node[T])
	n.data = data
	n.next = nil
	if l.tail == nil {
		l.head = n
		l.tail = n
		return
	}
	l.tail.next = n
	l.tail = n
}

func (l *List[T]) Insert(index int, data T) {
	l.length++
	n := l.pool.Get().(*node[T])
	n.data = data
	n.next = nil
	if index == 0 {
		n.next = l.head
		l.head = n
		return
	}
	if index == l.Len() {
		l.tail.next = n
		l.tail = n
		return
	}
	i := 0
	for p := l.head; p != nil; p = p.next {
		if i == index-1 {
			n.next = p.next
			p.next = n
			return
		}
		i++
	}
}

func (l *List[T]) Remove(index int) (rs T, ok bool) {
	l.length--
	if index == 0 {
		if l.head == nil {
			return
		}
		rs = l.head.data
		l.pool.Put(l.head)
		l.head = l.head.next
		ok = true
		return
	}
	i := 0
	for p := l.head; p != nil; p = p.next {
		if i == index-1 {
			if p.next == nil {
				return
			}
			rs = p.next.data
			l.pool.Put(p.next)
			p.next = p.next.next
			ok = true
			return
		}
		i++
	}
	return
}

func (l *List[T]) At(index int) (rs T, ok bool) {
	i := 0
	for p := l.head; p != nil; p = p.next {
		if i == index {
			rs = p.data
			ok = true
			return
		}
		i++
	}
	return
}

func (l *List[T]) Iterate(f func(T) bool) {
	for p := l.head; p != nil; p = p.next {
		if !f(p.data) {
			return
		}
	}
}

func (l *List[T]) Len() (rs int) {
	return l.length
}
