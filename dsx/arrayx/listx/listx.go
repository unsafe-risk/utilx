package listx

import "github.com/unsafe-risk/utilx/dsx"

var _ = (dsx.List[int])(&List[int]{})

type List[T any] []T

func (l *List[T]) Append(data T) {
	*l = append(*l, data)
}

func (l *List[T]) Insert(index int, data T) {
	*l = append(*l, data)
	copy((*l)[index+1:], (*l)[index:])
	(*l)[index] = data
}

func (l *List[T]) Remove(index int) (rs T, ok bool) {
	if index < 0 || index >= len(*l) {
		return
	}
	rs = (*l)[index]
	ok = true
	*l = append((*l)[:index], (*l)[index+1:]...)
	return
}

func (l *List[T]) At(index int) (rs T, ok bool) {
	if index < 0 || index >= len(*l) {
		return
	}
	rs = (*l)[index]
	ok = true
	return
}

func (l *List[T]) Iterate(f func(T) bool) {
	for _, v := range *l {
		if !f(v) {
			break
		}
	}
}

func (l *List[T]) Len() int {
	return len(*l)
}

func (l *List[T]) Cap() int {
	return cap(*l)
}
