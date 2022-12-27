package slicex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	for _, test := range []struct {
		slices []any
	}{
		{[]any{nil}},
		{[]any{1, 2, 3}},
		{[]any{1, "2", 3.2, []int{1, 2}, []string{"a", "b"}}},
	} {
		ori := New(test.slices...)
		cpy := ori.Copy()
		require.Equal(t, ori, cpy)
	}
}

func TestSetGet(t *testing.T) {
	for _, test := range []struct {
		sliceLen int
		idxs     []int
		vals     []any
	}{
		{1, []int{0}, []any{1}},
		{3, []int{0, 1, 2}, []any{1, 2, 3}},
		{5, []int{0, 1, 2, 3, 4}, []any{1, "2", 3.2, []int{1, 2}, []string{"a", "b"}}},

		// index exception
		{2, []int{-1, 6}, []any{1, 2}},
	} {
		slice := New[any]()
		slice.Extend(test.sliceLen)
		for i, idx := range test.idxs {
			slice.Set(idx, test.vals[i])
		}
		for i, idx := range test.idxs {
			require.Equal(t, test.vals[i], slice.Get(idx))
		}
	}
}

func TestSlice(t *testing.T) {
	for _, test := range []struct {
		ori        []any
		sliceStart int
		sliceEnd   int
		expect     []any
	}{
		{[]any{0, 1, 2}, 0, 0, []any{}},
		{[]any{0, 1, 2}, 0, 1, []any{0}},
		{[]any{0, 1, 2}, 0, 2, []any{0, 1}},
		{[]any{0, 1, 2}, 0, 3, []any{0, 1, 2}},
		{[]any{0, 1, 2}, 1, 3, []any{1, 2}},
		{[]any{0, 1, 2}, 2, 3, []any{2}},
		{[]any{0, 1, 2}, 3, 3, []any{}},

		// minus index
		{[]any{0, 1, 2}, 0, -1, []any{0, 1}},
		{[]any{0, 1, 2}, 0, -2, []any{0}},

		// index exception
		{[]any{0, 1, 2}, 0, 100, []any{0, 1, 2}},
		{[]any{0, 1, 2}, -100, 3, []any{0, 1, 2}},
		{[]any{0, 1, 2}, -100, -1, []any{0, 1}},
		{[]any{0, 1, 2}, -100, -3, []any{}},
	} {
		ori := New(test.ori...).Slice(test.sliceStart, test.sliceEnd)
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
	}
}

func TestSplice(t *testing.T) {
	for _, test := range []struct {
		ori         []any
		spliceStart int
		spliceEnd   int
		spliceVals  []any
		expect      []any
	}{
		{[]any{0, 1, 2}, 0, 0, []any{}, []any{0, 1, 2}},
		{[]any{0, 1, 2}, 0, 1, []any{3}, []any{3, 1, 2}},
		{[]any{0, 1, 2}, 0, 2, []any{3, 4}, []any{3, 4, 2}},
		{[]any{0, 1, 2}, 0, 3, []any{3, 4, 5}, []any{3, 4, 5}},
		{[]any{0, 1, 2}, 1, 1, []any{3}, []any{0, 3, 1, 2}},
		{[]any{0, 1, 2}, 1, 2, []any{3, 4}, []any{0, 3, 4, 2}},
		{[]any{0, 1, 2}, 1, 3, []any{3, 4, 5}, []any{0, 3, 4, 5}},
		{[]any{0, 1, 2}, 2, 3, []any{4, 5}, []any{0, 1, 4, 5}},
		{[]any{0, 1, 2}, 3, 3, []any{4}, []any{0, 1, 2, 4}},

		// index exception
		{[]any{0, 1, 2}, 0, 100, []any{3, 4, 5}, []any{3, 4, 5}},
		{[]any{0, 1, 2}, -100, 3, []any{3, 4, 5}, []any{3, 4, 5}},
	} {
		ori := New(test.ori...)
		ori.Splice(test.spliceStart, test.spliceEnd, test.spliceVals...)
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
	}
}

// TODO
/*
func TestInsert(t *testing.T) {
	for _, test := range []struct {
		ori    []any
		idx    int
		vals   []any
		expect []any
	}{
		{[]any{1, 2, 3}, 0, []any{0}, []any{0, 1, 2, 3}},
		{[]any{1, 2, 3}, 1, []any{0}, []any{1, 0, 2, 3}},
		{[]any{1, 2, 3}, 2, []any{0}, []any{1, 2, 0, 3}},
		{[]any{1, 2, 3}, 3, []any{0}, []any{1, 2, 3, 0}},
		{[]any{1, 2, 3}, 0, []any{4, 5, 6}, []any{4, 5, 6, 1, 2, 3}},
		{[]any{1, 2, 3}, 1, []any{4, 5, 6}, []any{1, 4, 5, 6, 2, 3}},
		{[]any{1, 2, 3}, 2, []any{4, 5, 6}, []any{1, 2, 4, 5, 6, 3}},
		{[]any{1, 2, 3}, 3, []any{4, 5, 6}, []any{1, 2, 3, 4, 5, 6}},

		{[]any{1, 2, 3}, 3, nil, []any{1, 2, 3}},
		{nil, 3, []any{4, 5, 6}, []any{1, 2, 3}},

		// index exception
		{[]any{1, 2, 3}, -100, []any{0}, []any{0, 1, 2, 3}},
		{[]any{1, 2, 3}, 100, []any{0}, []any{1, 2, 3, 0}},
	} {
		ori := New(test.ori...).Insert(test.idx, test.vals...)
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
	}
}

func TestDelete(t *testing.T) {
	slice := New[any](1, 2, 3, "4", 5.5, []byte{6, 6, 6}, []string{"7", "7", "7"}, 8.8, -9)
	fmt.Println(slice)
	slice.Delete(-1)
	fmt.Println(slice)
	slice.Delete(2)
	fmt.Println(slice)
	slice.Delete(1000)
	fmt.Println(slice)
}

func TestPushPop(t *testing.T) {
	ori := New[any](1, 2, 3).Push(4).Push(5.5).Push("안녕").UnShift("hello")

	fmt.Println(ori)
	fmt.Println("pop :", ori.Pop())
	fmt.Println("shift :", ori.Shift())
	fmt.Println("pop :", ori.Pop())
	fmt.Println("shift :", ori.Shift())
	fmt.Println("pop :", ori.Pop())
	fmt.Println("shift :", ori.Shift())
	fmt.Println("pop :", ori.Pop())
	fmt.Println("shift :", ori.Shift())
	fmt.Println("pop :", ori.Pop())
	fmt.Println("shift :", ori.Shift())
	fmt.Println("pop :", ori.Pop())
}

func TestRotate(t *testing.T) {
	ori := New[any](1, "2", 3, 4, 5, 6, 7, 8, 9)

	fmt.Println(ori.Rotate(0))
	fmt.Println(ori.Rotate(9))
	fmt.Println(ori.Rotate(3))
	fmt.Println(ori.Rotate(-3))
	fmt.Println(ori.Rotate(-3))
	fmt.Println(ori.Shuffle())
	fmt.Println(ori.Reverse())
}

func TestSplit(t *testing.T) {
	ori := New[any](1, "2", 3, 4, 5, 6, 7, 8, 9)
	split := ori.Split(1, 2, 3, 4)
	fmt.Println(split)
	split = ori.Split(1, 1, 1, 0, 0, 0, -3)
	fmt.Println(split)
}

func TestBatch(t *testing.T) {
	ori := New[any](1, "2", 3, 4, 5, 6, 7, 8, 9)

	batch := ori.Batch(-1)
	fmt.Println(batch)

	batch = ori.Batch(1)
	fmt.Println(batch)

	batch = ori.Batch(3)
	fmt.Println(batch)

	batch = ori.Batch(4)
	fmt.Println(batch)

	batch = ori.Batch(5)
	fmt.Println(batch)

	batch = ori.Batch(10000000000000000)
	fmt.Println(batch)
}

func TestCap(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	a1 := a[0:1]

	a = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	a3 := a[2:4]

	a = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	a9 := a[0:9]

	a = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	a11 := a[0:1:1]

	a = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	a33 := a[0:3:3]

	a = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	a35 := a[0:3:5]

	fmt.Println(a1, cap(a1))
	fmt.Println(a3, cap(a3))
	fmt.Println(a9, cap(a9))
	fmt.Println(a11, cap(a11))
	fmt.Println(a33, cap(a33))
	fmt.Println(a35, cap(a35))
}
*/
