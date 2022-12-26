package convx

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

func integerToInteger[T, R constraints.Integer | constraints.Float](i T) R {
	return R(i)
}

func intToString[T constraints.Signed](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

func uintToString[T constraints.Unsigned](i T) string {
	return strconv.FormatUint(uint64(i), 10)
}

func floatToString[T constraints.Float](i T) string {
	return strconv.FormatFloat(float64(i), 'f', -1, 64)
}

func stringToInteger[T constraints.Integer](s string) T {
	i, _ := strconv.ParseInt(s, 10, 64)
	return T(i)
}

func stringToUinteger[T constraints.Unsigned](s string) T {
	i, _ := strconv.ParseUint(s, 10, 64)
	return T(i)
}

func stringToFloat[T constraints.Float](s string) T {
	i, _ := strconv.ParseFloat(s, 64)
	return T(i)
}

func boolToString(b bool) string {
	return strconv.FormatBool(b)
}

func stringToBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func boolToInteger[T constraints.Integer](b bool) T {
	if b {
		return T(1)
	}
	return T(0)
}

func intToBool[T constraints.Integer](i T) bool {
	return i != T(0)
}

func boolToFloat[T constraints.Float](b bool) T {
	if b {
		return T(1)
	}
	return T(0)
}

func floatToBool[T constraints.Float](f T) bool {
	return f != T(0)
}

func stringToString(s string) string {
	return s
}

func boolToBool(b bool) bool {
	return b
}

func stringToBytes(s string) []byte {
	return []byte(s)
}

func bytesToString(b []byte) string {
	return string(b)
}

func bytesToBytes(b []byte) []byte {
	return b
}

func stringToRunes(s string) []rune {
	return []rune(s)
}

func runesToString(r []rune) string {
	return string(r)
}

func runesToRunes(r []rune) []rune {
	return r
}
