package convx

import (
	"reflect"
	"sync"
)

var converters = make(map[reflect.Type]map[reflect.Type]any)
var convertersLock = sync.RWMutex{}

func Register[T, R any](converter func(T) R) {
	convertersLock.Lock()
	defer convertersLock.Unlock()
	if converters[reflect.TypeOf(*new(T))] == nil {
		converters[reflect.TypeOf(*new(T))] = make(map[reflect.Type]any)
	}
	converters[reflect.TypeOf(*new(T))][reflect.TypeOf(*new(R))] = converter
}

func Into[T, R any](value T) (rs R, ok bool) {
	convertersLock.RLock()
	defer convertersLock.RUnlock()
	m := converters[reflect.TypeOf(value)]
	if m == nil {
		return
	}
	v := m[reflect.TypeOf(rs)]
	if v == nil {
		return
	}
	f, ok := v.(func(T) R)
	if !ok {
		return
	}
	rs = f(value)
	ok = true
	return
}

func IntoOr[T, R any](value T, or R) R {
	rs, ok := Into[T, R](value)
	if !ok {
		return or
	}
	return rs
}
