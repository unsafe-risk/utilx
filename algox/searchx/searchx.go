package searchx

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func Search[T constraints.Ordered](s []T, x T) (int, bool) {
	return SearchCustom(s, x, func(i, j T) int {
		if i < j {
			return -1
		}
		if i > j {
			return 1
		}
		return 0
	})
}

func SearchCustom[T any](s []T, x T, less func(i, j T) int) (int, bool) {
	return slices.BinarySearchFunc(s, x, less)
}
