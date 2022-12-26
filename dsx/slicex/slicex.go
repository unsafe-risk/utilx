package slicex

import (
	"fmt"
	"math/rand"
)

type Slice[T any] []T

func New[T any](v ...T) Slice[T] {
	return v
}

func (s Slice[T]) Copy() Slice[T] {
	return append(Slice[T]{}, s...)
}

func (s Slice[T]) Replace(i int, v T) Slice[T] {
	s[i] = v
	return s
}

func (s Slice[T]) Push(v ...T) Slice[T] {
	s = append(s, v...)
	return s
}

func (s Slice[T]) Pop() T {
	var v T
	if len(s) == 0 {
		return v
	}
	v = s[len(s)-1]
	s = s[:len(s)-1]
	return v
}

func (s Slice[T]) UnShift(v ...T) Slice[T] {
	return append(v, s...)
}

func (s Slice[T]) Shift() T {
	var v T
	if len(s) == 0 {
		return v
	}
	v = s[0]
	s = s[1:]
	return v
}

func (s Slice[T]) Slice(i, j int) Slice[T] {
	s = s[i:j]
	return s
}

func (s Slice[T]) Insert(i int, v ...T) Slice[T] {
	return s.Splice(i, i, v...)
}

func (s Slice[T]) Splice(i, j int, v ...T) Slice[T] {
	s = append(append(s[:i], v...), s[j:]...)
	return s
}

func (s Slice[T]) Delete(i int) Slice[T] {
	return s.Cut(i, i+1)
}

func (s Slice[T]) DeleteReplace(i int) Slice[T] {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (s Slice[T]) Cut(i, j int) Slice[T] {
	s = append(s[:i], s[j:]...)
	return s
}

func (s Slice[T]) Reverse() Slice[T] {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func (s Slice[T]) Rotate() Slice[T] {
	var r T = s[0]
	for i := 1; i < len(s); i++ {
		s[i-1] = s[i]
	}
	s[len(s)-1] = r
	return s
}

func (s Slice[T]) DeRotate() Slice[T] {
	var r T = s[len(s)-1]
	for i := 0; i < len(s)-1; i++ {
		s[len(s)-1-i] = s[len(s)-2-i]
	}
	s[len(s)-1] = r
	return s
}

func (s Slice[T]) Shuffle() Slice[T] {
	rand.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
	return s
}

func (s Slice[T]) Filter(f func(T) bool) Slice[T] {
	var n int = 0
	for _, v := range s {
		if f(v) {
			s[n] = v
			n++
		}
	}
	s = s[:n]
	return s
}

func (s Slice[T]) Batch(size int) []Slice[T] {
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

func (s Slice[T]) Length() int {
	return len(s)
}

func (s Slice[T]) Clear() {
	s = s[:0]
}
