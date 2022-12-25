package intersectionx

import (
	"errors"
	"sort"
)

type iTuple struct {
	Offset uint64
	Type   int8
}

type iSlice []iTuple

func (s iSlice) Len() int {
	return len(s)
}

func (s iSlice) Less(i, j int) bool {
	return s[i].Offset < s[j].Offset
}

func (s iSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s iSlice) Sort() {
	sort.Sort(s)
}

var ErrInvalidInterval = errors.New("invalid interval")

func intersection(s iSlice) (lower, upper uint64, ok bool) {
	var M uint64 = uint64(s.Len()) / 3
	var F uint64

	for F < M/2 {
		var End, Mid uint64

		for i := range s {
			End = End - uint64(s[i].Type)
			if End >= M-F {
				lower = s[i].Offset
				break
			}

			if s[i].Type == 0 {
				Mid = Mid + 1
			}
		}

		End = 0

		for i := len(s) - 1; i >= 0; i-- {
			End = End + uint64(s[i].Type)
			if End >= M-F {
				upper = s[i].Offset
				break
			}

			if s[i].Type == 0 {
				Mid = Mid + 1
			}
		}

		if lower <= upper && Mid <= F {
			return lower, upper, true
		}

		F = F + 1
		if F >= M/2 {
			break
		}
	}

	ok = false
	return
}

type Interval struct {
	Offset         uint64
	ConfidenceBand uint64
}

func Intersection(intervals ...Interval) (upper, lower uint64, ok bool) {
	is := make(iSlice, 0, len(intervals)*3)
	for i := range intervals {
		is = append(
			is,
			iTuple{intervals[i].Offset - intervals[i].ConfidenceBand, -1},
			iTuple{intervals[i].Offset, 0},
			iTuple{intervals[i].Offset + intervals[i].ConfidenceBand, 1},
		)
	}
	sort.Sort(&is)
	return intersection(is)
}
