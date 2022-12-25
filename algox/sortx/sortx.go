package sortx

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func Sort[T constraints.Ordered](s []T, reverse bool) {
	SortCustom(s, func(i, j T) bool {
		return i < j
	})
}

func SortReverse[T constraints.Ordered](s []T) {
	SortCustom(s, func(i, j T) bool {
		return i > j
	})
}

func SortStable[T constraints.Ordered](s []T) {
	SortStableCustom(s, func(i, j T) bool {
		return i < j
	})
}

func SortStableReverse[T constraints.Ordered](s []T) {
	SortStableCustom(s, func(i, j T) bool {
		return i > j
	})
}

func IsSorted[T constraints.Ordered](s []T) bool {
	return IsSortedCustom(s, func(i, j T) bool {
		return i < j
	})
}

func IsSortedReverse[T constraints.Ordered](s []T) bool {
	return IsSortedCustom(s, func(i, j T) bool {
		return i > j
	})
}

func SortCustom[T any](s []T, less func(i, j T) bool) {
	slices.SortFunc(s, less)
}

func SortStableCustom[T any](s []T, less func(i, j T) bool) {
	slices.SortStableFunc(s, less)
}

func IsSortedCustom[T any](s []T, less func(i, j T) bool) bool {
	return slices.IsSortedFunc(s, less)
}
