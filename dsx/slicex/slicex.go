package slicex

import (
	"fmt"
	"math/rand"
)

type Slice[T any] []T

func New[T any](v ...T) *Slice[T] {
	var s Slice[T] = v
	return &s
}

func (s *Slice[T]) Copy() *Slice[T] {
	copy := append(Slice[T]{}, *s...)
	return &copy
}

func (s *Slice[T]) Clear() {
	*s = (*s)[:0]
}

func (s *Slice[T]) Replace(i int, v T) *Slice[T] {
	if i < 0 {
		i = 0
	}
	if i >= len(*s) {
		i = len(*s) - 1
	}
	(*s)[i] = v
	return s
}

func (s *Slice[T]) Push(v ...T) *Slice[T] {
	*s = append(*s, v...)
	return s
}

func (s *Slice[T]) Pop() T {
	var r T
	if len(*s) == 0 {
		return r
	}
	r = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return r
}

func (s *Slice[T]) UnShift(v ...T) *Slice[T] {
	*s = append(v, *s...)
	return s
}

func (s *Slice[T]) Shift() T {
	var r T
	if len(*s) == 0 {
		return r
	}
	r = (*s)[0]
	*s = (*s)[1:]
	return r
}

func (s *Slice[T]) Slice(i, j int) *Slice[T] {
	if j <= 0 {
		j = len(*s) + j
	}
	i, j = preprocIndexException(i, j, len(*s))
	*s = (*s)[i:j]
	return s
}

func (s *Slice[T]) Insert(i int, v ...T) *Slice[T] {
	i, j := preprocIndexException(i, i, len(*s))
	*s = append(append((*s)[:i], v...), (*s)[j:]...)
	return s
}

func (s *Slice[T]) Splice(i, j int, v ...T) *Slice[T] {
	if i < 0 {
		i = 0
	}
	if j > len(*s) {
		j = len(*s)
	}
	*s = append(append((*s)[:i], v...), (*s)[j:]...)
	return s
}

func (s *Slice[T]) Delete(i int) *Slice[T] {
	if i < 0 || i >= len(*s) {
		return s
	}
	*s = append((*s)[:i], (*s)[i+1:]...)
	return s
}

func (s *Slice[T]) DeleteReplace(i int) *Slice[T] {
	if i < 0 || i >= len(*s) {
		return s
	}
	(*s)[i] = (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return s
}

func (s *Slice[T]) Cut(i, j int) *Slice[T] {
	i, j = preprocIndexException(i, j, len(*s))
	*s = append((*s)[:i], (*s)[j:]...)
	return s
}

func (s *Slice[T]) Reverse() *Slice[T] {
	for i, j := 0, len(*s)-1; i < j; i, j = i+1, j-1 {
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	}
	return s
}

func (s *Slice[T]) Rotate(n int) *Slice[T] {
	if n > 0 {
		index := n % len(*s)
		*s = append((*s)[index:], (*s)[:index]...)
	} else if n < 0 {
		index := -n % len(*s)
		*s = append((*s)[len(*s)-index:], (*s)[:len(*s)-index]...)
	}
	return s
}

func (s *Slice[T]) Shuffle() *Slice[T] {
	rand.Shuffle(len(*s), func(i, j int) {
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	})
	return s
}

func (s *Slice[T]) Filter(f func(T) bool) *Slice[T] {
	n := 0
	for _, v := range *s {
		if f(v) {
			(*s)[n] = v
			n++
		}
	}
	*s = (*s)[:n]
	return s
}

func (s Slice[T]) Batch(size int) []Slice[T] {
	if size <= 0 {
		size = 1
	}
	batches := make([]Slice[T], 0, (len(s)+size-1)/size)
	for size < len(s) {
		s, batches = s[size:], append(batches, s[0:size:size])
	}
	return append(batches, s)
}

func (s Slice[T]) Join(sep string) string {
	var r string
	for i, v := range s {
		if i != 0 {
			r += sep
		}
		r += fmt.Sprintf("%v", v)
	}
	return r
}

func (s Slice[T]) IndexOf(f func(T) bool) int {
	for i, v := range s {
		if f(v) {
			return i
		}
	}
	return -1
}

func (s Slice[T]) Contains(f func(T) bool) bool {
	for _, v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

func (s Slice[T]) Each(f func(T)) {
	for _, v := range s {
		f(v)
	}
}

func (s Slice[T]) ForEach(f func(int, T)) {
	for i, v := range s {
		f(i, v)
	}
}

func (s Slice[T]) Len() int {
	return len(s)
}

// 인덱스 전처리
func preprocIndexException(i, j, len int) (int, int) {
	if i < 0 {
		i = 0
	}
	if j > len {
		j = len
	}
	if i > j {
		j = i
	}
	return i, j
}

func (s Slice[T]) Finalize() []T {
	return s
}
