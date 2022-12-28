package queuex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnqueueDequeue(t *testing.T) {
	for _, test := range []struct {
		node   []any
		expect []any
	}{
		{[]any{1}, []any{1}},
		{[]any{1, 2}, []any{1, 2}},
		{[]any{1, 2, 3}, []any{1, 2, 3}},
	} {
		ori := New[any](len(test.node))
		for _, v := range test.node {
			ori.Enqueue(v)
		}

		var idx int = 0
		for {
			r, ok := ori.Dequeue()
			if ok != true {
				break
			}
			require.Equal(t, test.expect[idx], r)
			idx++
		}
	}
}
