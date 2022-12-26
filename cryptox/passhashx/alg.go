package passhashx

import (
	"github.com/unsafe-risk/utilx/cryptox/passhashx/internal"
	"golang.org/x/crypto/argon2"
)

func alg_Argon2ID_High(password []byte, salt []byte) []byte {
	const Argon2ID_Time = 20
	const Argon2ID_Memory = 64 * 1024
	const Argon2ID_Parallelism = 4
	return argon2.IDKey(
		password, salt,
		Argon2ID_Time, Argon2ID_Memory, Argon2ID_Parallelism,
		hash_len(internal.Parameter_Argon2ID_High),
	)
}

func alg_Argon2ID_Low(password []byte, salt []byte) []byte {
	const Argon2ID_Time = 32
	const Argon2ID_Memory = 4 * 1024
	const Argon2ID_Parallelism = 1
	return argon2.IDKey(
		password, salt,
		Argon2ID_Time, Argon2ID_Memory, Argon2ID_Parallelism,
		hash_len(internal.Parameter_Argon2ID_Low),
	)
}

func alg_Argon2ID_Mobile_High(password []byte, salt []byte) []byte {
	const Argon2ID_Time = 4
	const Argon2ID_Memory = 37 * 1024
	const Argon2ID_Parallelism = 1
	return argon2.IDKey(
		password, salt,
		Argon2ID_Time, Argon2ID_Memory, Argon2ID_Parallelism,
		hash_len(internal.Parameter_Argon2ID_Mobile_High),
	)
}

func alg_Argon2ID_Mobile_Low(password []byte, salt []byte) []byte {
	const Argon2ID_Time = 8
	const Argon2ID_Memory = 15 * 1024
	const Argon2ID_Parallelism = 1
	return argon2.IDKey(
		password, salt,
		Argon2ID_Time, Argon2ID_Memory, Argon2ID_Parallelism,
		hash_len(internal.Parameter_Argon2ID_Mobile_Low),
	)
}
