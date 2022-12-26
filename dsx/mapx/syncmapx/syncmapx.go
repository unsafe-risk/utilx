package syncmapx

import "sync"

// goroutine-safe map
type Map[K comparable, V any] struct {
	m   map[K]V
	mtx sync.RWMutex
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m: make(map[K]V),
	}
}

func (m *Map[K, V]) Get(k K) V {
	m.mtx.RLock()
	val := m.m[k]
	m.mtx.RUnlock()
	return val
}

func (m *Map[K, V]) Has(k K) bool {
	m.mtx.RLock()
	_, ok := m.m[k]
	m.mtx.RUnlock()
	return ok
}

func (m *Map[K, V]) Size() int {
	m.mtx.RLock()
	size := len(m.m)
	m.mtx.RUnlock()
	return size
}

func (t *Map[K, V]) Keys() []K {
	t.mtx.RLock()
	keys := make([]K, 0, len(t.m))
	for k := range t.m {
		keys = append(keys, k)
	}
	t.mtx.RUnlock()
	return keys
}

func (t *Map[K, V]) Values() []V {
	t.mtx.RLock()
	items := make([]V, 0, len(t.m))
	for _, v := range t.m {
		items = append(items, v)
	}
	t.mtx.RUnlock()
	return items
}

func (m *Map[K, V]) Set(k K, v V) {
	m.mtx.Lock()
	m.m[k] = v
	m.mtx.Unlock()
}

func (m *Map[K, V]) Delete(k K) {
	m.mtx.Lock()
	delete(m.m, k)
	m.mtx.Unlock()
}

func (m *Map[K, V]) Clear() {
	m.mtx.Lock()
	m.m = make(map[K]V)
	m.mtx.Unlock()
}
