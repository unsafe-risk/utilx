package sklmapx

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

const MAX_HEIGHT = 16

func (l *skl[T, V]) rand_height() int {
	var h int = 1
	for ; h < MAX_HEIGHT; h++ {
		if splitmix64(&l.seed)&0x1 == 1 {
			return h
		}
	}
	return h
}

type skl[T constraints.Ordered, V any] struct {
	head *sknode[T, V]
	tail *sknode[T, V]
	seed uint64
	len  uint
}

type sknode[T constraints.Ordered, V any] struct {
	key    T
	value  V
	next   [MAX_HEIGHT]*sknode[T, V]
	height uint16
}

func (s *sknode[T, V]) Height() int {
	return int(s.height)
}

func (l *skl[T, V]) search(key T, ctx *[MAX_HEIGHT]*sknode[T, V]) *sknode[T, V] {
	x := l.head
	for i := MAX_HEIGHT - 1; i >= 0; i-- {
		for x.next[i] != l.tail && x.next[i].key < key {
			x = x.next[i]
		}
		ctx[i] = x
	}
	x = x.next[0]
	if x == l.tail || x.key != key {
		return nil
	}
	return x
}

func (l *skl[T, V]) insert(key T, value V, ctx *[MAX_HEIGHT]*sknode[T, V]) *sknode[T, V] {
	x := l.search(key, ctx)
	if x != nil {
		x.value = value
		return x
	}

	level := l.rand_height()
	node := &sknode[T, V]{key: key, value: value, height: uint16(level)}

	for i := 0; i < level; i++ {
		node.next[i] = ctx[i].next[i]
		ctx[i].next[i] = node
	}
	l.len++

	return node
}

func (l *skl[T, V]) delete(key T, ctx *[MAX_HEIGHT]*sknode[T, V]) bool {
	x := l.search(key, ctx)
	if x == nil {
		return false
	}

	for i := 0; i < x.Height(); i++ {
		ctx[i].next[i] = x.next[i]
	}
	l.len--

	return true
}

type SkipListMap[T constraints.Ordered, V any] struct {
	skl skl[T, V]
}

func (m *SkipListMap[T, V]) Set(key T, value V) {
	var ctx [MAX_HEIGHT]*sknode[T, V]
	m.skl.insert(key, value, &ctx)
}

func (m *SkipListMap[T, V]) Del(key T) bool {
	var ctx [MAX_HEIGHT]*sknode[T, V]
	return m.skl.delete(key, &ctx)
}

func (m *SkipListMap[T, V]) Get(key T) (V, bool) {
	var ctx [MAX_HEIGHT]*sknode[T, V]
	x := m.skl.search(key, &ctx)
	if x == nil {
		var zero V
		return zero, false
	}
	return x.value, true
}

func New[T constraints.Ordered, V any]() *SkipListMap[T, V] {
	head := &sknode[T, V]{}
	tail := &sknode[T, V]{}
	for i := 0; i < MAX_HEIGHT; i++ {
		head.next[i] = tail
	}
	seed := root_next()
	return &SkipListMap[T, V]{skl: skl[T, V]{head: head, tail: tail, seed: seed}}
}

func (m *SkipListMap[T, V]) Iterator() *Iterator[T, V] {
	return &Iterator[T, V]{list: m}
}

func (m *SkipListMap[T, V]) String() string {
	var buf []byte
	i := Iterator[T, V]{list: m}
	buf = append(buf, "sklmapx["...)
	var idx int
	for {
		key, value, ok := i.Next()
		if !ok {
			break
		}
		if idx > 0 {
			buf = append(buf, ' ')
		}
		idx++
		buf = fmt.Appendf(buf, "%v:%v", key, value)
	}
	buf = append(buf, ']')
	return string(buf)
}

func (m *SkipListMap[T, V]) Keys() []T {
	var keys []T = make([]T, 0, m.Len())
	i := Iterator[T, V]{list: m}
	for {
		key, _, ok := i.Next()
		if !ok {
			break
		}
		keys = append(keys, key)
	}
	return keys
}

func (m *SkipListMap[T, V]) Values() []V {
	var values []V = make([]V, 0, m.Len())
	i := Iterator[T, V]{list: m}
	for {
		_, value, ok := i.Next()
		if !ok {
			break
		}
		values = append(values, value)
	}
	return values
}

func (m *SkipListMap[T, V]) Len() int {
	return int(m.skl.len)
}

func (m *SkipListMap[T, V]) ToMap() map[T]V {
	mm := make(map[T]V, m.Len())
	i := Iterator[T, V]{list: m}
	for {
		key, value, ok := i.Next()
		if !ok {
			break
		}
		mm[key] = value
	}
	return mm
}

func (m *SkipListMap[T, V]) Clear() {
	// Clear the map by creating a new one.
	*m = *New[T, V]()
}

type Iterator[T constraints.Ordered, V any] struct {
	list *SkipListMap[T, V]
	ctx  [MAX_HEIGHT]*sknode[T, V]
	node *sknode[T, V]
}

func (i *Iterator[T, V]) Seek(key T) (T, V, bool) {
	i.node = i.list.skl.search(key, &i.ctx)
	if i.node == nil {
		var zero T
		var zeroV V
		return zero, zeroV, false
	}
	return i.node.key, i.node.value, true
}

func (i *Iterator[T, V]) Rewind() {
	i.node = nil
}

func (i *Iterator[T, V]) Next() (T, V, bool) {
	if i.node == nil {
		i.node = i.list.skl.head.next[0]
	} else {
		i.node = i.node.next[0]
	}
	if i.node == i.list.skl.tail {
		var zero T
		var zeroV V
		return zero, zeroV, false
	}
	return i.node.key, i.node.value, true
}
