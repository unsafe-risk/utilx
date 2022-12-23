package passhashx

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"

	"github.com/unsafe-risk/utilx/cryptox/passhashx/internal"
)

const PASSHASH_VERSION = 1

type SecurityLevel uint16

const (
	SecurityLevelHigh SecurityLevel = 1
	SecurityLevelLow  SecurityLevel = 2
)

func salt_len(param internal.Parameter) uint8 {
	switch param {
	case internal.Parameter_Argon2ID_High:
		return 16
	case internal.Parameter_Argon2ID_Low:
		return 16
	}

	return 16
}

func hash_len(param internal.Parameter) uint32 {
	switch param {
	case internal.Parameter_Argon2ID_High:
		return 32
	case internal.Parameter_Argon2ID_Low:
		return 16
	}

	return 32
}

var ErrReadSaltFailed = errors.New("read salt failed")
var ErrInvalidSecurityLevel = errors.New("invalid security level")

func Hash(password []byte, sl SecurityLevel) ([]byte, error) {
	var salt []byte

	switch sl {
	case SecurityLevelHigh:
		const param = internal.Parameter_Argon2ID_High
		salt = make([]byte, salt_len(param))
	case SecurityLevelLow:
		const param = internal.Parameter_Argon2ID_Low
		salt = make([]byte, salt_len(param))
	default:
		return nil, ErrInvalidSecurityLevel
	}

	n, err := rand.Read(salt)
	if err != nil || n != len(salt) {
		return nil, ErrReadSaltFailed
	}

	var hash []byte
	switch sl {
	case SecurityLevelHigh:
		hash = alg_Argon2ID_High(password, salt)
	case SecurityLevelLow:
		hash = alg_Argon2ID_Low(password, salt)
	}

	switch sl {
	case SecurityLevelHigh:
		const param = internal.Parameter_Argon2ID_High
		return internal.New_PasswordHash(
			param,
			salt,
			hash,
		), err
	case SecurityLevelLow:
		const param = internal.Parameter_Argon2ID_Low
		return internal.New_PasswordHash(
			param,
			salt,
			hash,
		), err
	}
	return nil, ErrInvalidSecurityLevel
}

func HashBase64(password []byte, sl SecurityLevel) (string, error) {
	hash, err := Hash(password, sl)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(hash), nil
}

var ErrInvalidPassHash = errors.New("invalid passhash")
var ErrUnknownAlgorithm = errors.New("unknown algorithm")
var ErrHashMismatch = errors.New("hash mismatch")

func Verify(password []byte, phash []byte) error {
	data := internal.PasswordHash(phash)
	if !data.Vstruct_Validate() {
		return ErrInvalidPassHash
	}

	var hash []byte

	switch data.Param() {
	case internal.Parameter_Argon2ID_High:
		hash = alg_Argon2ID_High(password, data.Salt())
	case internal.Parameter_Argon2ID_Low:
		hash = alg_Argon2ID_Low(password, data.Salt())
	default:
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
