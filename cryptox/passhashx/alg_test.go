package passhashx

import (
	"testing"
)

var password = []byte("0123456789abcdef0123456789abcdef")
var salt = []byte("0123456789abcdef0123456789abcdef")

func BenchmarkAlgArgon2ID_High(b *testing.B) {
	for i := 0; i < b.N; i++ {
		alg_Argon2ID_High(password, salt)
	}
}

func BenchmarkAlgArgon2ID_Low(b *testing.B) {
	for i := 0; i < b.N; i++ {
		alg_Argon2ID_Low(password, salt)
	}
}
