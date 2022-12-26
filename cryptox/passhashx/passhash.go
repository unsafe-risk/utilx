package passhashx

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"

	"github.com/unsafe-risk/utilx/cryptox/passhashx/internal"
)

var (
	ErrInvalidPassHash      = errors.New("invalid passhash")
	ErrUnknownAlgorithm     = errors.New("unknown algorithm")
	ErrHashMismatch         = errors.New("hash mismatch")
	ErrReadSaltFailed       = errors.New("read salt failed")
	ErrInvalidSecurityLevel = errors.New("invalid security level")
)

const PASSHASH_VERSION = 1

type SecurityLevel int16

const (
	SecurityLevelHigh       SecurityLevel = 1
	SecurityLevelLow        SecurityLevel = 2
	SecurityLevelMobileHigh SecurityLevel = 3
	SecurityLevelMobileLow  SecurityLevel = 4
)

const internal_Parameter_MAX = 4

var salt_len_arr = [internal_Parameter_MAX]uint8{
	internal.Parameter_Argon2ID_High:        16,
	internal.Parameter_Argon2ID_Low:         16,
	internal.Parameter_Argon2ID_Mobile_High: 16,
	internal.Parameter_Argon2ID_Mobile_Low:  16,
}

func salt_len(param internal.Parameter) uint8 {
	if param < internal.Parameter(len(salt_len_arr)) {
		return salt_len_arr[param]
	}
	return 16
}

var hash_len_arr = [internal_Parameter_MAX]uint32{
	internal.Parameter_Argon2ID_High:        32,
	internal.Parameter_Argon2ID_Low:         32,
	internal.Parameter_Argon2ID_Mobile_High: 32,
	internal.Parameter_Argon2ID_Mobile_Low:  32,
}

func hash_len(param internal.Parameter) uint32 {
	if param < internal.Parameter(len(hash_len_arr)) {
		return hash_len_arr[param]
	}
	return 32
}

var alg_arr = [internal_Parameter_MAX]func(password []byte, salt []byte) []byte{
	internal.Parameter_Argon2ID_High:        alg_Argon2ID_High,
	internal.Parameter_Argon2ID_Low:         alg_Argon2ID_Low,
	internal.Parameter_Argon2ID_Mobile_High: alg_Argon2ID_Mobile_High,
	internal.Parameter_Argon2ID_Mobile_Low:  alg_Argon2ID_Mobile_Low,
}

func security_level(sl SecurityLevel) (internal.Parameter, bool) {
	switch sl {
	case SecurityLevelHigh:
		return internal.Parameter_Argon2ID_High, true
	case SecurityLevelLow:
		return internal.Parameter_Argon2ID_Low, true
	case SecurityLevelMobileHigh:
		return internal.Parameter_Argon2ID_Mobile_High, true
	case SecurityLevelMobileLow:
		return internal.Parameter_Argon2ID_Mobile_Low, true
	}
	return internal.Parameter_Argon2ID_High, false
}

func Hash(password []byte, sl SecurityLevel) ([]byte, error) {
	var salt []byte
	var ok bool
	var param internal.Parameter
	if param, ok = security_level(sl); !ok {
		return nil, ErrInvalidSecurityLevel
	}

	salt = make([]byte, salt_len(param))
	n, err := rand.Read(salt)
	if err != nil || n != len(salt) {
		return nil, ErrReadSaltFailed
	}

	var hash []byte
	if param < internal.Parameter(len(alg_arr)) {
		hash = alg_arr[param](password, salt)
	} else {
		return nil, ErrUnknownAlgorithm
	}

	return internal.New_PasswordHash(
		param,
		salt,
		hash,
	), nil
}

func HashBase64(password []byte, sl SecurityLevel) (string, error) {
	hash, err := Hash(password, sl)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(hash), nil
}

func Verify(password []byte, phash []byte) error {
	data := internal.PasswordHash(phash)
	if !data.Vstruct_Validate() {
		return ErrInvalidPassHash
	}

	var hash []byte

	if data.Param() < internal.Parameter(len(alg_arr)) {
		alg := alg_arr[data.Param()]
		if alg != nil {
			hash = alg(password, data.Salt())
		} else {
			return ErrUnknownAlgorithm
		}
	} else {
		return ErrUnknownAlgorithm
	}

	if subtle.ConstantTimeCompare(hash, data.Hash()) != 1 {
		return ErrHashMismatch
	}

	return nil
}

func VerifyBase64(password []byte, phash string) error {
	hash, err := base64.RawURLEncoding.DecodeString(phash)
	if err != nil {
		return err
	}
	return Verify(password, hash)
}
