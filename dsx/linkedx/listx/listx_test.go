package listx_test

import (
	"testing"

	"github.com/unsafe-risk/utilx/dsx/linkedx/listx"
)

func TestListX(t *testing.T) {
	l := listx.New[int]()
	for i := 0; i < 100; i++ {
		l.Append(i)
	}
	for i := 0; i < 10; i++ {
		l.Insert(i, i)
	}
	for i := 0; i < 10; i++ {
		if v, ok := l.At(i); !ok || v != i {
			t.Fatalf("listx: failed to insert")
		}
	}
	for i := 0; i < 10; i++ {
		if v, ok := l.At(i + 10); !ok || v != i {
			t.Fatalf("listx: failed to insert")
		}
	}
	for i := 19; i >= 0; i -= 2 {
		l.Remove(i)
	}
	is := []int{0, 2, 4, 6, 8, 0, 2, 4, 6, 8}
	for i, p := range is {
		if v, ok := l.At(i); !ok || v != p {
			t.Fatalf("listx: failed to remove")
		}
	}
}
