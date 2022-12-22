package stackx

type Stack[T any] struct {
	data []T
}

func New[T any]() *Stack[T] {
	return &Stack[T]{data: make([]T, 0)}
}

func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
}

func (s *Stack[T]) Pop() (rs T, ok bool) {
	if len(s.data) == 0 {
		return
	}
	rs = s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return rs, true
}

func (s *Stack[T]) Peek() (rs T, ok bool) {
	if len(s.data) == 0 {
		return
	}
	rs = s.data[len(s.data)-1]
	return rs, true
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.data) == 0
}
