package bitsetx

import (
	"golang.org/x/exp/constraints"
)

type BitSet[T constraints.Unsigned] struct {
	data T
}

func (b *BitSet[T]) Set(i ...int) {
	for _, i := range i {
		b.data |= (1 << i)
	}
}

func (b *BitSet[T]) Clear(i ...int) {
	for _, i := range i {
		b.data &= ^(1 << i)
	}
}

func (b *BitSet[T]) Toggle(i ...int) {
	for _, i := range i {
		b.data ^= (1 << i)
	}
}

func (b *BitSet[T]) Test(i ...int) bool {
	t := false
	for _, i := range i {
		t = t || b.data&(1<<i) != 0
	}
	return t
}

func (b *BitSet[T]) IsEmpty() bool {
	return b.data == 0
}

func (b *BitSet[T]) IsFull() bool {
	return b.data == ^T(0)
}

func (b *BitSet[T]) ClearAll() {
	b.data = 0
}

func (b *BitSet[T]) ToggleAll() {
	b.data = ^b.data
}

func (b *BitSet[T]) SetAll() {
	b.data = ^T(0)
}
