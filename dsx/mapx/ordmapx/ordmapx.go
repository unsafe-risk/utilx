package ordmapx

import (
	"github.com/unsafe-risk/utilx/dsx/mapx/sklmapx"
	"golang.org/x/exp/constraints"
)

type OrdMap[K constraints.Ordered, V any] struct {
	m   map[K]V
	skl *sklmapx.SkipListMap[K, struct{}]
}

func New[K constraints.Ordered, V any]() *OrdMap[K, V] {
	return &OrdMap[K, V]{
		m:   make(map[K]V),
		skl: sklmapx.New[K, struct{}](),
	}
}

func (m *OrdMap[K, V]) Len() int {
	return m.skl.Len()
}

func (m *OrdMap[K, V]) Get(key K) (V, bool) {
	v, ok := m.m[key]
	return v, ok
}

func (m *OrdMap[K, V]) Set(key K, value V) {
	if _, ok := m.m[key]; !ok {
		m.skl.Set(key, struct{}{})
	}
	m.m[key] = value
}

func (m *OrdMap[K, V]) Del(key K) {
	if _, ok := m.m[key]; ok {
		m.skl.Del(key)
		delete(m.m, key)
	}
}

func (m *OrdMap[K, V]) Clear() {
	m.m = make(map[K]V)
	m.skl.Clear()
}

func (m *OrdMap[K, V]) Keys() []K {
	return m.skl.Keys()
}

func (m *OrdMap[K, V]) Values() []V {
	iter := m.skl.Iterator()
	values := make([]V, 0, m.skl.Len())
	for {
		k, _, ok := iter.Next()
		if !ok {
			break
		}
		values = append(values, m.m[k])
	}
	return values
}

func (m *OrdMap[K, V]) Range(f func(key K, value V) (stop bool)) {
	iter := m.skl.Iterator()
	for {
		k, _, ok := iter.Next()
		if !ok {
			break
		}
		if stop := f(k, m.m[k]); stop {
			break
		}
	}
}
