package ordmapx

// TODO - use btree to order
type Map[K comparable, V any] struct {
	m map[K]V
	// order list[K]
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m: make(map[K]V),
	}
}

func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	value, ok = m.m[key]
	return
}

/*
// set comparable
func (m *Map[K, V]) Set(key K, value V) {
	if _, ok := m.m[key]; !ok {
		m.order = append(m.order, key)
	}
	m.m[key] = value
}

func (m *Map[K, V]) Delete(key K) {
	delete(m.m, key)
	for i, k := range m.order {
		if k == key {
			m.order = append(m.order[:i], m.order[i+1:]...)
			return
		}
	}
}

func (m *Map[K, V]) Len() int {
	return len(m.m)
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	for _, key := range m.order {
		if !f(key, m.m[key]) {
			return
		}
	}
}

func (m *Map[K, V]) Keys() []K {
	return m.order
}

func (m *Map[K, V]) Values() []V {
	values := make([]V, len(m.order))
	for i, key := range m.order {
		values[i] = m.m[key]
	}
	return values
}
*/
