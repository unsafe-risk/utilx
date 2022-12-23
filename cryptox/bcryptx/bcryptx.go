package bcryptx

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// 길이 제한
// 비밀번호가 72 바이트 초과일 경우 73 bt 이후가 달라도 같다고 판정하기 때문
const (
	DefMaxLen = 72
)

var (
	ErrPasswordTooLong = fmt.Errorf("pw len exceed | pw len must be less than %v", DefMaxLen)
)

func Encrypt(pw []byte) ([]byte, error) {
	if len(pw) > DefMaxLen {
		return nil, ErrPasswordTooLong
	}
	hash, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func Decrypt(hash, pw []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, pw)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, err
		}
		return false, err
	}
	return true, nil
}
