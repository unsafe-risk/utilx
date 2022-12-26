package slicex

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	fn_test := func(s ...any) {
		ori := New[any](s)
		cpy := ori.Copy()
		fmt.Println(cpy.Join(", "))
		require.Equal(t, ori, cpy)
	}

	fn_test()
	fn_test(nil)
	fn_test(1, 2, 3)
	fn_test(1, "2", 3.2, []int{1, 2}, []string{"a", "b"})
}

func TestPushPop(t *testing.T) {
	ori := New[any](1, "2", 3).Push(1).Push(2.2).Push("안녕").UnShift("hello")

	fmt.Println(ori.Pop())
	fmt.Println(ori.Shift())
	fmt.Println(ori.Pop())
	fmt.Println(ori.Shift())
	fmt.Println(ori.Pop())
	fmt.Println(ori.Shift())
	fmt.Println(ori.Pop())
	fmt.Println(ori.Shift())
	fmt.Println(ori.Pop())
	fmt.Println(ori.Shift())
	fmt.Println(ori.Pop())
}

func TestSlice(t *testing.T) {
	slice := New[any](1, 2, 3, "4", 5.5, []byte{6, 6, 6}, []string{"7", "7", "7"}, 8.8, -9)
	fmt.Println(slice.Slice(0, 0))
	fmt.Println(slice.Slice(1, 8))
	fmt.Println(slice.Slice(1, -2))
	fmt.Println(slice.Splice(1, 3))
}

func TestSplice(t *testing.T) {
	ori := New[any](1, "2", 3, 4, 5, 6, 7, 8, 9)

	fmt.Println(ori.Splice(-1, 100, 100, 99, 98))
	fmt.Println(ori.Splice(2, 4, 97, 96, 95))
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
