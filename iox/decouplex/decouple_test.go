package decouplex_test

import (
	"sync/atomic"
	"testing"

	"github.com/unsafe-risk/utilx/iox/decouplex"
)

func TestDrop1(t *testing.T) {
	var dropped uint64
	dropfn := func(b []byte) (int, error) {
		atomic.AddUint64(&dropped, 1)
		return len(b), nil
	}
	d := decouplex.NewDecoupler(1, dropfn)
	d.Write([]byte("hello"))
	d.Write([]byte("world"))
	if dropped != 1 {
		t.Fatalf("expected 1 dropped, got %d", dropped)
	}
}

func drain(d *decouplex.Decoupler) {
	for d.Next(nil) {
	}
}

func TestDrop2(t *testing.T) {
	var dropped uint64
	dropfn := func(b []byte) (int, error) {
		atomic.AddUint64(&dropped, 1)
		return len(b), nil
	}
	d := decouplex.NewDecoupler(2, dropfn)

	d.Write([]byte("hello"))
	d.Write([]byte("world"))
	drain(d)
	d.Write([]byte("world"))
	if dropped != 0 {
		t.Fatalf("expected 0 dropped, got %d", dropped)
	}

	drain(d)
	dropped = 0
	d.Write([]byte("hello"))
	d.Write([]byte("world"))
	d.Write([]byte("world"))
	d.Write([]byte("world"))
	if dropped != 2 {
		t.Fatalf("expected 2 dropped, got %d", dropped)
	}

	drain(d)
	dropped = 0
	d.Write([]byte("hello"))
	d.Write([]byte("world"))
	if dropped != 0 {
		t.Fatalf("expected 0 dropped, got %d", dropped)
	}
}
