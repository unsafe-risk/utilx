package dequex_test

import (
	"testing"

	"github.com/unsafe-risk/utilx/dsx/dequex"
)

func TestDeque(t *testing.T) {
	d := dequex.New[int]()
	for i := 0; i < 1000; i++ {
		d.PushFront(i)
		d.PushBack(i)
	}

	for i := 0; i < 1000; i++ {
		v, ok := d.PopFront()
		if !ok {
			t.Fatal("PopFront failed")
		}
		if v != 999-i {
			t.Fatal("PopFront failed")
		}
		v, ok = d.PopBack()
		if !ok {
			t.Fatal("PopBack failed")
		}
		if v != 999-i {
			t.Fatal("PopBack failed")
		}
	}
}
