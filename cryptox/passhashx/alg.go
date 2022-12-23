package passhashx

import (
	"github.com/unsafe-risk/utilx/cryptox/passhashx/internal"
	"golang.org/x/crypto/argon2"
)

func alg_Argon2ID_High(password []byte, salt []byte) []byte {
	const Argon2ID_High_Time = 20
	const Argon2ID_High_Memory = 64 * 1024
	const Argon2ID_High_Parallelism = 4
	return argon2.IDKey(
		password, salt,
		Argon2ID_High_Time, Argon2ID_High_Memory, Argon2ID_High_Parallelism,
		hash_len(internal.Parameter_Argon2ID_High),
	)
}

func alg_Argon2ID_Low(password []byte, salt []byte) []byte {
	const Argon2ID_Low_Time = 32
	const Argon2ID_Low_Memory = 4 * 1024
	const Argon2ID_Low_Parallelism = 1
	return argon2.IDKey(
		password, salt,
		Argon2ID_Low_Time, Argon2ID_Low_Memory, Argon2ID_Low_Parallelism,
		hash_len(internal.Parameter_Argon2ID_Low),
	)
}
