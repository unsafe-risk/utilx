package slicex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	for _, test := range []struct {
		slices []any
	}{
		{[]any{}},
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
		expects  []any
	}{
		{1, []int{0}, []any{1}, []any{1}},
		{3, []int{0, 1, 2}, []any{1, 2, 3}, []any{1, 2, 3}},
		{5, []int{0, 1, 2, 3, 4}, []any{1, "2", 3.2, []int{1, 2}, []string{"a", "b"}}, []any{1, "2", 3.2, []int{1, 2}, []string{"a", "b"}}},
		{10, []int{0, 2, 4, 6, 8}, []any{1, 2, 3, 4, 5}, []any{1, nil, 2, nil, 3, nil, 4, nil, 5, nil}},

		// index exception
		{4, []int{-1, 6}, []any{1, 2}, []any{1, nil, nil, 2}},
	} {
		ori := New[any]()
		ori.Extend(test.sliceLen)
		for i, idx := range test.idxs {
			ori.Set(idx, test.vals[i])
		}
		for i, idx := range test.idxs {
			require.Equal(t, test.vals[i], ori.Get(idx))
		}

		exp := New(test.expects...)
		require.Equal(t, exp, ori)
	}
}

func TestSlice(t *testing.T) {
	for _, test := range []struct {
		ori      []any
		idxStart int
		idxEnd   int
		expect   []any
	}{
		{[]any{0, 1, 2}, 0, 0, []any{}},
		{[]any{0, 1, 2}, 0, 1, []any{0}},
		{[]any{0, 1, 2}, 0, 2, []any{0, 1}},
		{[]any{0, 1, 2}, 0, 3, []any{0, 1, 2}},
		{[]any{0, 1, 2}, 1, 3, []any{1, 2}},
		{[]any{0, 1, 2}, 2, 3, []any{2}},
		{[]any{0, 1, 2}, 3, 3, []any{}},

		{[]any{}, 0, 1, []any{}},

		// minus index
		{[]any{0, 1, 2}, 0, -1, []any{0, 1}},
		{[]any{0, 1, 2}, 0, -2, []any{0}},

		// index exception
		{[]any{0, 1, 2}, 0, 100, []any{0, 1, 2}},
		{[]any{0, 1, 2}, -100, 3, []any{0, 1, 2}},
		{[]any{0, 1, 2}, -100, -1, []any{0, 1}},
		{[]any{0, 1, 2}, -100, -3, []any{}},
	} {
		ori := New(test.ori...).Slice(test.idxStart, test.idxEnd)
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
	}
}

func TestSplice(t *testing.T) {
	for _, test := range []struct {
		ori      []any
		idxStart int
		idxEnd   int
		vals     []any
		expect   []any
	}{
		{[]any{0, 1, 2}, 0, 0, nil, []any{0, 1, 2}},
		{[]any{0, 1, 2}, 0, 0, []any{}, []any{0, 1, 2}},
		{[]any{0, 1, 2}, 0, 1, []any{3}, []any{3, 1, 2}},
		{[]any{0, 1, 2}, 0, 2, []any{3, 4}, []any{3, 4, 2}},
		{[]any{0, 1, 2}, 0, 3, []any{3, 4, 5}, []any{3, 4, 5}},
		{[]any{0, 1, 2}, 1, 1, []any{3}, []any{0, 3, 1, 2}},
		{[]any{0, 1, 2}, 1, 2, []any{3, 4}, []any{0, 3, 4, 2}},
		{[]any{0, 1, 2}, 1, 3, []any{3, 4, 5}, []any{0, 3, 4, 5}},
		{[]any{0, 1, 2}, 2, 3, []any{4, 5}, []any{0, 1, 4, 5}},
		{[]any{0, 1, 2}, 3, 3, []any{4}, []any{0, 1, 2, 4}},

		{[]any{1, 2, 3}, 0, 1, []any{}, []any{2, 3}},
		{[]any{}, 0, 1, []any{3, 4, 5}, []any{3, 4, 5}},
		{[]any{}, 0, 1, []any{}, []any{}},

		// index exception
		{[]any{0, 1, 2}, 0, 100, []any{3, 4, 5}, []any{3, 4, 5}},
		{[]any{0, 1, 2}, -100, 3, []any{3, 4, 5}, []any{3, 4, 5}},
	} {
		ori := New(test.ori...)
		ori.Splice(test.idxStart, test.idxEnd, test.vals...)
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
	}
}

func TestExpand(t *testing.T) {
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
		{nil, 3, []any{4, 5, 6}, []any{4, 5, 6}},

		// index exception
		{[]any{1, 2, 3}, -100, []any{0}, []any{0, 1, 2, 3}},
		{[]any{1, 2, 3}, 100, []any{0}, []any{1, 2, 3, 0}},
	} {
		ori := New(test.ori...).Expand(test.idx, test.vals...)
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
	}
}

func TestInsert(t *testing.T) {
	for _, test := range []struct {
		ori    []any
		idx    int
		val    any
		expect []any
	}{
		{[]any{1, 2, 3}, 0, any(0), []any{0, 1, 2, 3}},
		{[]any{1, 2, 3}, 1, any(0), []any{1, 0, 2, 3}},
		{[]any{1, 2, 3}, 2, any(0), []any{1, 2, 0, 3}},
		{[]any{1, 2, 3}, 3, any(0), []any{1, 2, 3, 0}},

		{[]any{1, 2, 3}, 3, any(nil), []any{1, 2, 3, nil}},

		// index exception
		{[]any{1, 2, 3}, -100, any(0), []any{0, 1, 2, 3}},
		{[]any{1, 2, 3}, 100, any(0), []any{1, 2, 3, 0}},
	} {
		ori := New(test.ori...).Insert(test.idx, test.val)
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
	}
}

func TestCut(t *testing.T) {
	for _, test := range []struct {
		ori      []any
		idxStart int
		idxEnd   int
		expect   []any
	}{
		{[]any{1, 2, 3}, 0, 0, []any{1, 2, 3}},
		{[]any{1, 2, 3}, 0, 1, []any{2, 3}},
		{[]any{1, 2, 3}, 0, 2, []any{3}},
		{[]any{1, 2, 3}, 0, 3, []any{}},
		{[]any{1, 2, 3}, 1, 1, []any{1, 2, 3}},
		{[]any{1, 2, 3}, 1, 2, []any{1, 3}},
		{[]any{1, 2, 3}, 1, 3, []any{1}},
		{[]any{1, 2, 3}, 2, 2, []any{1, 2, 3}},
		{[]any{1, 2, 3}, 2, 3, []any{1, 2}},
		{[]any{1, 2, 3}, 3, 3, []any{1, 2, 3}},

		{[]any{}, 0, 1, []any{}},

		// index exception
		{[]any{1, 2, 3}, 0, 100, []any{}},
		{[]any{1, 2, 3}, -100, 3, []any{}},
	} {
		ori := New(test.ori...).Cut(test.idxStart, test.idxEnd)
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
	}
}

func TestDelete(t *testing.T) {
	for _, test := range []struct {
		ori    []any
		idx    int
		expect []any
	}{
		{[]any{1, 2, 3}, 0, []any{2, 3}},
		{[]any{1, 2, 3}, 1, []any{1, 3}},
		{[]any{1, 2, 3}, 2, []any{1, 2}},

		{[]any{}, 1, []any{}},

		// index exception
		{[]any{1, 2, 3}, -100, []any{1, 2, 3}},
		{[]any{1, 2, 3}, 100, []any{1, 2, 3}},
	} {
		ori := New(test.ori...).Delete(test.idx)
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
	}
}

func TestPush(t *testing.T) {
	for _, test := range []struct {
		ori    []any
		val    any
		expect []any
	}{
		{[]any{1, 2, 3}, 4, []any{1, 2, 3, 4}},
		{[]any{}, 4, []any{4}},
		{nil, 4, []any{4}},
	} {
		ori := New(test.ori...).Push(test.val)
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
	}
}

func TestPop(t *testing.T) {
	for _, test := range []struct {
		ori       []any
		expect    []any
		expectPop any
	}{
		{[]any{1, 2, 3}, []any{1, 2}, 3},
		{[]any{1}, []any{}, 1},
		{[]any{}, []any{}, nil},
		{nil, nil, any(nil)},
	} {
		ori := New(test.ori...)
		pop := ori.Pop()
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
		require.Equal(t, test.expectPop, pop)
	}
}

func TestUnshift(t *testing.T) {
	for _, test := range []struct {
		ori    []any
		val    any
		expect []any
	}{
		{[]any{1, 2, 3}, 4, []any{4, 1, 2, 3}},
		{[]any{}, 4, []any{4}},
		{nil, 4, []any{4}},
	} {
		ori := New(test.ori...).UnShift(test.val)
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
	}
}

func TestShift(t *testing.T) {
	for _, test := range []struct {
		ori         []any
		expect      []any
		expectShift any
	}{
		{[]any{1, 2, 3}, []any{2, 3}, 1},
		{[]any{1}, []any{}, 1},
		{[]any{}, []any{}, nil},
		{nil, nil, any(nil)},
	} {
		ori := New(test.ori...)
		pop := ori.Shift()
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
		require.Equal(t, test.expectShift, pop)
	}
}

func TestReverse(t *testing.T) {
	for _, test := range []struct {
		ori    []any
		expect []any
	}{
		{[]any{1, 2, 3}, []any{3, 2, 1}},
		{[]any{1, 2, 3, 4}, []any{4, 3, 2, 1}},

		{[]any{1}, []any{1}},
		{[]any{}, []any{}},
	} {
		ori := New(test.ori...).Reverse()
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
	}
}

func TestRotate(t *testing.T) {
	for _, test := range []struct {
		ori    []any
		rotate int
		expect []any
	}{
		{[]any{1, 2, 3}, 0, []any{1, 2, 3}},
		{[]any{1, 2, 3}, 1, []any{2, 3, 1}},
		{[]any{1, 2, 3}, 2, []any{3, 1, 2}},
		{[]any{1, 2, 3}, 3, []any{1, 2, 3}},
		{[]any{1, 2, 3}, 4, []any{2, 3, 1}},

		{[]any{1, 2, 3}, -1, []any{3, 1, 2}},
		{[]any{1, 2, 3}, -2, []any{2, 3, 1}},
		{[]any{1, 2, 3}, -3, []any{1, 2, 3}},
		{[]any{1, 2, 3}, -4, []any{3, 1, 2}},

		{[]any{1}, 1, []any{1}},
		{[]any{}, 1, []any{}},
	} {
		ori := New(test.ori...).Rotate(test.rotate)
		exp := New(test.expect...)
		require.Equal(t, exp, ori)
	}
}

func TestSplit(t *testing.T) {
	for _, test := range []struct {
		ori    []any
		idxs   []int
		expect [][]any
	}{
		{[]any{1, 2, 3}, []int{1}, [][]any{{1}, {2, 3}}},
		{[]any{1, 2, 3}, []int{2}, [][]any{{1, 2}, {3}}},
		{[]any{1, 2, 3}, []int{3}, [][]any{{1, 2, 3}}},

		{[]any{1, 2, 3}, []int{}, [][]any{{1, 2, 3}}},
		{[]any{1, 2, 3}, []int{0}, [][]any{{}, {1, 2, 3}}},
		{[]any{1, 2, 3}, []int{0, 1}, [][]any{{}, {1}, {2, 3}}},

		{[]any{1, 2, 3}, []int{1, 1}, [][]any{{1}, {2}, {3}}},
		{[]any{1, 2, 3}, []int{1, 2}, [][]any{{1}, {2, 3}}},
		{[]any{1, 2, 3}, []int{2}, [][]any{{1, 2}, {3}}},
		{[]any{1, 2, 3}, []int{3}, [][]any{{1, 2, 3}}},

		{[]any{1, 2, 3}, []int{1, 1, 1}, [][]any{{1}, {2}, {3}}},

		// minus index
		{[]any{1, 2, 3}, []int{1, -1}, [][]any{{1}, {1}, {1, 2, 3}}},
		{[]any{1, 2, 3}, []int{2, -2}, [][]any{{1, 2}, {1, 2}, {1, 2, 3}}},
		{[]any{1, 2, 3}, []int{3, -3}, [][]any{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}}},

		// index exception
		{[]any{1, 2, 3}, []int{-100}, [][]any{{}, {1, 2, 3}}},
		{[]any{1, 2, 3}, []int{100, -100}, [][]any{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}}},
		{[]any{1, 2, 3}, []int{-100, 100, -100}, [][]any{{}, {1, 2, 3}, {1, 2, 3}, {1, 2, 3}}},
	} {
		oris := New(test.ori...).Split(test.idxs...)
		exps := make([]Slice[any], len(test.expect))
		for i, exp := range test.expect {
			exps[i] = exp
		}

		require.Equal(t, exps, oris)
	}
}

func TestBatch(t *testing.T) {
	for _, test := range []struct {
		ori    []any
		size   int
		expect [][]any
	}{
		{[]any{1, 2, 3, 4, 5}, 1, [][]any{{1}, {2}, {3}, {4}, {5}}},
		{[]any{1, 2, 3, 4, 5}, 2, [][]any{{1, 2}, {3, 4}, {5}}},
		{[]any{1, 2, 3, 4, 5}, 3, [][]any{{1, 2, 3}, {4, 5}}},
		{[]any{1, 2, 3, 4, 5}, 4, [][]any{{1, 2, 3, 4}, {5}}},
		{[]any{1, 2, 3, 4, 5}, 5, [][]any{{1, 2, 3, 4, 5}}},

		// index exception
		{[]any{1, 2, 3, 4, 5}, -1, [][]any{{1}, {2}, {3}, {4}, {5}}},
		{[]any{1, 2, 3, 4, 5}, 0, [][]any{{1}, {2}, {3}, {4}, {5}}},
		{[]any{1, 2, 3, 4, 5}, 100, [][]any{{1, 2, 3, 4, 5}}},
	} {
		oris := New(test.ori...).Batch(test.size)
		exps := make([]Slice[any], len(test.expect))
		for i, exp := range test.expect {
			exps[i] = exp
		}
		require.Equal(t, exps, oris)
	}
}
