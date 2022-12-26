package sklmapx_test

import (
	"reflect"
	"testing"

	"github.com/unsafe-risk/utilx/dsx/mapx/sklmapx"
)

func TestSkipListMap(t *testing.T) {
	skip := sklmapx.New[string, string]()
	skip.Set("b", "b")
	skip.Set("f", "f")
	skip.Set("d", "d")
	skip.Set("c", "c")
	skip.Set("a", "a")
	skip.Set("e", "e")

	if v, ok := skip.Get("a"); !ok || v != "a" {
		t.Errorf("expected a, got %v", v)
	}

	if v, ok := skip.Get("b"); !ok || v != "b" {
		t.Errorf("expected b, got %v", v)
	}

	if v, ok := skip.Get("c"); !ok || v != "c" {
		t.Errorf("expected c, got %v", v)
	}

	if v, ok := skip.Get("d"); !ok || v != "d" {
		t.Errorf("expected d, got %v", v)
	}

	if v, ok := skip.Get("e"); !ok || v != "e" {
		t.Errorf("expected e, got %v", v)
	}

	if v, ok := skip.Get("f"); !ok || v != "f" {
		t.Errorf("expected f, got %v", v)
	}

	if v, ok := skip.Get("g"); ok {
		t.Errorf("expected no value, got %v", v)
	}

	skip.Del("a")

	if v, ok := skip.Get("a"); ok {
		t.Errorf("expected no value, got %v", v)
	}

	skip.Del("b")

	if v, ok := skip.Get("b"); ok {
		t.Errorf("expected no value, got %v", v)
	}

	skip.Del("g") // delete non-existent key

	if v, ok := skip.Get("g"); ok {
		t.Errorf("expected no value, got %v", v)
	}

	skip.Del("c")

	if v, ok := skip.Get("c"); ok {
		t.Errorf("expected no value, got %v", v)
	}

	// Remaining keys: d, e, f

	m := skip.ToMap()
	if len(m) != 3 {
		t.Errorf("expected 3 keys, got %v", len(m))
	}

	if v, ok := m["d"]; !ok || v != "d" {
		t.Errorf("expected d, got %v", v)
	}

	if v, ok := m["e"]; !ok || v != "e" {
		t.Errorf("expected e, got %v", v)
	}

	if v, ok := m["f"]; !ok || v != "f" {
		t.Errorf("expected f, got %v", v)
	}

	if v, ok := m["g"]; ok {
		t.Errorf("expected no value, got %v", v)
	}

	// Test iterator

	iter0 := skip.Iterator()
	if k, v, ok := iter0.Next(); !ok || k != "d" || v != "d" {
		t.Errorf("expected d, got %v", k)
	}

	if k, v, ok := iter0.Next(); !ok || k != "e" || v != "e" {
		t.Errorf("expected e, got %v", k)
	}

	iter1 := skip.Iterator()

	if k, v, ok := iter1.Next(); !ok || k != "d" || v != "d" {
		t.Errorf("expected d, got %v", k)
	}

	if k, v, ok := iter0.Next(); !ok || k != "f" || v != "f" {
		t.Errorf("expected f, got %v", k)
	}

	if k, v, ok := iter1.Next(); !ok || k != "e" || v != "e" {
		t.Errorf("expected e, got %v", k)
	}

	if k, _, ok := iter0.Next(); ok {
		t.Errorf("expected no value, got %v", k)
	}

	if k, v, ok := iter1.Next(); !ok || k != "f" || v != "f" {
		t.Errorf("expected f, got %v", k)
	}

	if k, _, ok := iter1.Next(); ok {
		t.Errorf("expected no value, got %v", k)
	}

	keys := skip.Keys()

	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %v", len(keys))
	}

	if !reflect.DeepEqual(keys, []string{"d", "e", "f"}) {
		t.Errorf("expected [d e f], got %v", keys)
	}

	values := skip.Values()

	if len(values) != 3 {
		t.Errorf("expected 3 values, got %v", len(values))
	}

	if !reflect.DeepEqual(values, []string{"d", "e", "f"}) {
		t.Errorf("expected [d e f], got %v", values)
	}

	skip.Clear()

	if v, ok := skip.Get("d"); ok {
		t.Errorf("expected no value, got %v", v)
	}

	if v, ok := skip.Get("e"); ok {
		t.Errorf("expected no value, got %v", v)
	}

	if v, ok := skip.Get("f"); ok {
		t.Errorf("expected no value, got %v", v)
	}

	if v, ok := skip.Get("g"); ok {
		t.Errorf("expected no value, got %v", v)
	}

	if len(skip.Keys()) != 0 {
		t.Errorf("expected 0 keys, got %v", len(skip.Keys()))
	}

	// Test Stringer

	if s := skip.String(); s != "sklmapx[]" {
		t.Errorf("expected sklmapx[], got %v", s)
	}

	skip.Set("abcde", "12345")
	skip.Set("abc", "123")
	skip.Set("ab", "12")
	skip.Set("a", "1")

	if s := skip.String(); s != "sklmapx[a:1 ab:12 abc:123 abcde:12345]" {
		t.Errorf("expected sklmapx[a:1 ab:12 abc:123 abcde:12345], got %v", s)
	}

	// Test Seek

	iter := skip.Iterator()
	iter.Seek("ab")

	if k, v, ok := iter.Next(); !ok || k != "abc" || v != "123" {
		t.Errorf("expected abc, got %v", k)
	}

	if k, v, ok := iter.Next(); !ok || k != "abcde" || v != "12345" {
		t.Errorf("expected abcde, got %v", k)
	}

	iter.Seek("a")

	if k, v, ok := iter.Next(); !ok || k != "ab" || v != "12" {
		t.Errorf("expected ab, got %v", k)
	}

	iter.Rewind()

	m = skip.ToMap()
	keys = skip.Keys()
	var idx int

	for k, v, ok := iter.Next(); ok; {
		if m[k] != v {
			t.Errorf("expected %v, got %v", m[k], v)
		}

		if keys[idx] != k {
			t.Errorf("expected %v, got %v", keys[idx], k)
		}

		idx++
		k, v, ok = iter.Next()
	}

	if idx != len(m) {
		t.Errorf("expected %v, got %v", len(m), idx)
	}
}
