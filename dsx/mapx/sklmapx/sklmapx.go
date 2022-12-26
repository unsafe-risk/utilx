package sklmapx

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

const MAX_HEIGHT = 16

func (l *skl[K, V]) rand_height() int {
	var h int = 1
	for ; h < MAX_HEIGHT; h++ {
		if splitmix64(&l.seed)&0x1 == 1 {
			return h
		}
	}
	return h
}

type skl[K constraints.Ordered, V any] struct {
	head *sknode[K, V]
	tail *sknode[K, V]
	seed uint64
	len  uint
}

type sknode[K constraints.Ordered, V any] struct {
	key    K
	value  V
	next   [MAX_HEIGHT]*sknode[K, V]
	height uint16
}

func (s *sknode[K, V]) Height() int {
	return int(s.height)
}

func (l *skl[K, V]) search(key K, ctx *[MAX_HEIGHT]*sknode[K, V]) *sknode[K, V] {
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

func (l *skl[K, V]) insert(key K, value V, ctx *[MAX_HEIGHT]*sknode[K, V]) *sknode[K, V] {
	x := l.search(key, ctx)
	if x != nil {
		x.value = value
		return x
	}

	level := l.rand_height()
	node := &sknode[K, V]{key: key, value: value, height: uint16(level)}

	for i := 0; i < level; i++ {
		node.next[i] = ctx[i].next[i]
		ctx[i].next[i] = node
	}
	l.len++

	return node
}

func (l *skl[K, V]) delete(key K, ctx *[MAX_HEIGHT]*sknode[K, V]) bool {
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

type SkipListMap[K constraints.Ordered, V any] struct {
	skl skl[K, V]
}

func (m *SkipListMap[K, V]) Set(key K, value V) {
	var ctx [MAX_HEIGHT]*sknode[K, V]
	m.skl.insert(key, value, &ctx)
}

func (m *SkipListMap[K, V]) Del(key K) bool {
	var ctx [MAX_HEIGHT]*sknode[K, V]
	return m.skl.delete(key, &ctx)
}

func (m *SkipListMap[K, V]) Get(key K) (V, bool) {
	var ctx [MAX_HEIGHT]*sknode[K, V]
	x := m.skl.search(key, &ctx)
	if x == nil {
		var zero V
		return zero, false
	}
	return x.value, true
}

func New[K constraints.Ordered, V any]() *SkipListMap[K, V] {
	head := &sknode[K, V]{}
	tail := &sknode[K, V]{}
	for i := 0; i < MAX_HEIGHT; i++ {
		head.next[i] = tail
	}
	seed := root_next()
	return &SkipListMap[K, V]{skl: skl[K, V]{head: head, tail: tail, seed: seed}}
}

func (m *SkipListMap[K, V]) Iterator() *Iterator[K, V] {
	return &Iterator[K, V]{list: m}
}

func (m *SkipListMap[K, V]) String() string {
	var buf []byte
	i := Iterator[K, V]{list: m}
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

func (m *SkipListMap[K, V]) Keys() []K {
	var keys []K = make([]K, 0, m.Len())
	i := Iterator[K, V]{list: m}
	for {
		key, _, ok := i.Next()
		if !ok {
			break
		}
		keys = append(keys, key)
	}
	return keys
}

func (m *SkipListMap[K, V]) Values() []V {
	var values []V = make([]V, 0, m.Len())
	i := Iterator[K, V]{list: m}
	for {
		_, value, ok := i.Next()
		if !ok {
			break
		}
		values = append(values, value)
	}
	return values
}

func (m *SkipListMap[K, V]) Len() int {
	return int(m.skl.len)
}

func (m *SkipListMap[K, V]) ToMap() map[K]V {
	mm := make(map[K]V, m.Len())
	i := Iterator[K, V]{list: m}
	for {
		key, value, ok := i.Next()
		if !ok {
			break
		}
		mm[key] = value
	}
	return mm
}

func (m *SkipListMap[K, V]) Clear() {
	// Clear the map by creating a new one.
	*m = *New[K, V]()
}

type Iterator[K constraints.Ordered, V any] struct {
	list *SkipListMap[K, V]
	ctx  [MAX_HEIGHT]*sknode[K, V]
	node *sknode[K, V]
}

func (i *Iterator[K, V]) Seek(key K) (K, V, bool) {
	i.node = i.list.skl.search(key, &i.ctx)
	if i.node == nil {
		var zero K
		var zeroV V
		return zero, zeroV, false
	}
	return i.node.key, i.node.value, true
}

func (i *Iterator[K, V]) Rewind() {
	i.node = nil
}

func (i *Iterator[K, V]) Next() (K, V, bool) {
	if i.node == nil {
		i.node = i.list.skl.head.next[0]
	} else {
		i.node = i.node.next[0]
	}
	if i.node == i.list.skl.tail {
		var zero K
		var zeroV V
		return zero, zeroV, false
	}
	return i.node.key, i.node.value, true
}
