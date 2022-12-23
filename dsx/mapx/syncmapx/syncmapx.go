package syncmapx

import "sync"

// goroutine-safe map
type Map[T any] struct {
	m   map[string]T
	mtx sync.RWMutex
}

func New[T any]() *Map[T] {
	return &Map[T]{
		m: make(map[string]T),
	}
}

func (m *Map[T]) Get(key string) T {
	m.mtx.RLock()
	val := m.m[key]
	m.mtx.RUnlock()
	return val
}

func (m *Map[T]) Has(key string) bool {
	m.mtx.RLock()
	_, ok := m.m[key]
	m.mtx.RUnlock()
	return ok
}

func (m *Map[T]) Size() int {
	m.mtx.RLock()
	size := len(m.m)
	m.mtx.RUnlock()
	return size
}

func (t *Map[T]) Keys() []string {
	t.mtx.RLock()
	keys := make([]string, 0, len(t.m))
	for k := range t.m {
		keys = append(keys, k)
	}
	t.mtx.RUnlock()
	return keys
}

func (t *Map[T]) Values() []T {
	t.mtx.RLock()
	items := make([]T, 0, len(t.m))
	for _, v := range t.m {
		items = append(items, v)
	}
	t.mtx.RUnlock()
	return items
}

func (m *Map[T]) Set(key string, value T) {
	m.mtx.Lock()
	m.m[key] = value
	m.mtx.Unlock()
}

func (m *Map[T]) Delete(key string) {
	m.mtx.Lock()
	delete(m.m, key)
	m.mtx.Unlock()
}

func (m *Map[T]) Clear() {
	m.mtx.Lock()
	m.m = make(map[string]T)
	m.mtx.Unlock()
}
