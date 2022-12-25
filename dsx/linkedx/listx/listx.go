package listx

import "github.com/unsafe-risk/utilx/dsx"

type node[T any] struct {
	data T
	next *node[T]
}

var _ = (dsx.List[int])(&List[int]{})

type List[T any] struct {
	head *node[T]
	tail *node[T]
}

func New[T any]() *List[T] {
	return &List[T]{}
}

func (l *List[T]) Append(data T) {
	n := &node[T]{}
	n.data = data
	if l.tail == nil {
		l.head = n
		l.tail = n
		return
	}
	l.tail.next = n
	l.tail = n
}

func (l *List[T]) Insert(index int, data T) {
	n := &node[T]{}
	n.data = data
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
	if index == 0 {
		if l.head == nil {
			return
		}
		rs = l.head.data
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
	for p := l.head; p != nil; p = p.next {
		rs++
	}
	return
}
