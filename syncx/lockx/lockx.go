package lockx

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

func Lock[T any](l *T, f func(*T)) {
	p := l
	for !atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(l)), unsafe.Pointer(p), nil) {
		runtime.Gosched()
		p = l
	}
	f(p)
	for !atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(l)), nil, unsafe.Pointer(p)) {
		runtime.Gosched()
	}
}
