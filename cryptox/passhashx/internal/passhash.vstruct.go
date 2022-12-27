package internal

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unsafe"
)

type _ = strings.Builder
type _ = unsafe.Pointer

var _ = math.Float32frombits
var _ = math.Float64frombits
var _ = strconv.FormatInt
var _ = strconv.FormatUint
var _ = strconv.FormatFloat
var _ = fmt.Sprint

type Parameter uint8

const (
	Parameter_Argon2ID_High        Parameter = 0
	Parameter_Argon2ID_Low         Parameter = 1
	Parameter_Argon2ID_Mobile_High Parameter = 2
	Parameter_Argon2ID_Mobile_Low  Parameter = 3
	VSTRUCT_ENUM_Parameter_MAX               = 3
)

func (e Parameter) String() string {
	switch e {
	case Parameter_Argon2ID_High:
		return "Argon2ID_High"
	case Parameter_Argon2ID_Low:
		return "Argon2ID_Low"
	case Parameter_Argon2ID_Mobile_High:
		return "Argon2ID_Mobile_High"
	case Parameter_Argon2ID_Mobile_Low:
		return "Argon2ID_Mobile_Low"
	}
	return ""
}

func (e Parameter) Match(
	onArgon2ID_High func(),
	onArgon2ID_Low func(),
	onArgon2ID_Mobile_High func(),
	onArgon2ID_Mobile_Low func(),
) {
	switch e {
	case Parameter_Argon2ID_High:
		onArgon2ID_High()
	case Parameter_Argon2ID_Low:
		onArgon2ID_Low()
	case Parameter_Argon2ID_Mobile_High:
		onArgon2ID_Mobile_High()
	case Parameter_Argon2ID_Mobile_Low:
		onArgon2ID_Mobile_Low()
	}
}

func (e Parameter) MatchS(s struct {
	onArgon2ID_High        func()
	onArgon2ID_Low         func()
	onArgon2ID_Mobile_High func()
	onArgon2ID_Mobile_Low  func()
}) {
	switch e {
	case Parameter_Argon2ID_High:
		s.onArgon2ID_High()
	case Parameter_Argon2ID_Low:
		s.onArgon2ID_Low()
	case Parameter_Argon2ID_Mobile_High:
		s.onArgon2ID_Mobile_High()
	case Parameter_Argon2ID_Mobile_Low:
		s.onArgon2ID_Mobile_Low()
	}
}

var _ struct{} = func() struct{} {
	var _v [1]struct{}
	_v[Parameter_Argon2ID_High-0] = struct{}{}
	_v[Parameter_Argon2ID_Low-1] = struct{}{}
	_v[Parameter_Argon2ID_Mobile_High-2] = struct{}{}
	_v[Parameter_Argon2ID_Mobile_Low-3] = struct{}{}
	return struct{}{}
}()

type PasswordHash []byte

func (s PasswordHash) Param() Parameter {
	return Parameter(s[0])
}

func (s PasswordHash) Salt() []byte {
	_ = s[8]
	var __off0 uint64 = 17
	var __off1 uint64 = uint64(s[1]) |
		uint64(s[2])<<8 |
		uint64(s[3])<<16 |
		uint64(s[4])<<24 |
		uint64(s[5])<<32 |
		uint64(s[6])<<40 |
		uint64(s[7])<<48 |
		uint64(s[8])<<56
	return []byte(s[__off0:__off1])
}

func (s PasswordHash) Hash() []byte {
	_ = s[16]
	var __off0 uint64 = uint64(s[1]) |
		uint64(s[2])<<8 |
		uint64(s[3])<<16 |
		uint64(s[4])<<24 |
		uint64(s[5])<<32 |
		uint64(s[6])<<40 |
		uint64(s[7])<<48 |
		uint64(s[8])<<56
	var __off1 uint64 = uint64(s[9]) |
		uint64(s[10])<<8 |
		uint64(s[11])<<16 |
		uint64(s[12])<<24 |
		uint64(s[13])<<32 |
		uint64(s[14])<<40 |
		uint64(s[15])<<48 |
		uint64(s[16])<<56
	return []byte(s[__off0:__off1])
}

func (s PasswordHash) Vstruct_Validate() bool {
	if len(s) < 17 {
		return false
	}

	_ = s[16]

	var __off0 uint64 = 17
	var __off1 uint64 = uint64(s[1]) |
		uint64(s[2])<<8 |
		uint64(s[3])<<16 |
		uint64(s[4])<<24 |
		uint64(s[5])<<32 |
		uint64(s[6])<<40 |
		uint64(s[7])<<48 |
		uint64(s[8])<<56
	var __off2 uint64 = uint64(s[9]) |
		uint64(s[10])<<8 |
		uint64(s[11])<<16 |
		uint64(s[12])<<24 |
		uint64(s[13])<<32 |
		uint64(s[14])<<40 |
		uint64(s[15])<<48 |
		uint64(s[16])<<56
	var __off3 uint64 = uint64(len(s))
	return __off0 <= __off1 && __off1 <= __off2 && __off2 <= __off3
}

func (s PasswordHash) String() string {
	if !s.Vstruct_Validate() {
		return "PasswordHash (invalid)"
	}
	var __b strings.Builder
	__b.WriteString("PasswordHash {")
	__b.WriteString("Param: ")
	__b.WriteString(s.Param().String())
	__b.WriteString(", ")
	__b.WriteString("Salt: ")
	__b.WriteString(fmt.Sprint(s.Salt()))
	__b.WriteString(", ")
	__b.WriteString("Hash: ")
	__b.WriteString(fmt.Sprint(s.Hash()))
	__b.WriteString("}")
	return __b.String()
}

func Serialize_PasswordHash(dst PasswordHash, Param Parameter, Salt []byte, Hash []byte) PasswordHash {
	_ = dst[16]
	dst[0] = byte(Param)

	var __index = uint64(17)
	__tmp_1 := uint64(len(Salt)) + __index
	dst[1] = byte(__tmp_1)
	dst[2] = byte(__tmp_1 >> 8)
	dst[3] = byte(__tmp_1 >> 16)
	dst[4] = byte(__tmp_1 >> 24)
	dst[5] = byte(__tmp_1 >> 32)
	dst[6] = byte(__tmp_1 >> 40)
	dst[7] = byte(__tmp_1 >> 48)
	dst[8] = byte(__tmp_1 >> 56)
	copy(dst[__index:__tmp_1], Salt)
	__index += uint64(len(Salt))
	__tmp_2 := uint64(len(Hash)) + __index
	dst[9] = byte(__tmp_2)
	dst[10] = byte(__tmp_2 >> 8)
	dst[11] = byte(__tmp_2 >> 16)
	dst[12] = byte(__tmp_2 >> 24)
	dst[13] = byte(__tmp_2 >> 32)
	dst[14] = byte(__tmp_2 >> 40)
	dst[15] = byte(__tmp_2 >> 48)
	dst[16] = byte(__tmp_2 >> 56)
	copy(dst[__index:__tmp_2], Hash)
	return dst
}

func New_PasswordHash(Param Parameter, Salt []byte, Hash []byte) PasswordHash {
	var __vstruct__size = 17 + len(Salt) + len(Hash)
	var __vstruct__buf = make(PasswordHash, __vstruct__size)
	__vstruct__buf = Serialize_PasswordHash(__vstruct__buf, Param, Salt, Hash)
	return __vstruct__buf
}
