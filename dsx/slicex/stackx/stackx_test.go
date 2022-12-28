package stackx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPushPop(t *testing.T) {
	for _, test := range []struct {
		data   []any
		expect []any
	}{
		{[]any{1}, []any{1}},
		{[]any{1, 2}, []any{2, 1}},
		{[]any{1, 2, 3}, []any{3, 2, 1}},
	} {
		ori := New[any]()
		for _, v := range test.data {
			ori.Push(v)
		}

		var idx int = 0
		for {
			r, ok := ori.Pop()
			if ok != true {
				break
			}
			require.Equal(t, test.expect[idx], r)
			idx++
		}
	}
}
