package passhashx_test

import (
	"testing"

	"github.com/unsafe-risk/utilx/cryptox/passhashx"
)

func TestHashVerify(t *testing.T) {
	const password = "my super secret password"
	hash0, err := passhashx.Hash([]byte(password), passhashx.SecurityLevelHigh)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	if passhashx.Verify([]byte(password), hash0) != nil {
		t.Fatal("password verification failed")
		t.FailNow()
	}

	hash1, err := passhashx.Hash([]byte(password), passhashx.SecurityLevelLow)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	if passhashx.Verify([]byte(password), hash1) != nil {
		t.Fatal("password verification failed")
		t.FailNow()
	}

	hash2, err := passhashx.Hash([]byte(password), passhashx.SecurityLevelMobileHigh)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	if passhashx.Verify([]byte(password), hash2) != nil {
		t.Fatal("password verification failed")
		t.FailNow()
	}

	hash3, err := passhashx.Hash([]byte(password), passhashx.SecurityLevelMobileLow)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	if passhashx.Verify([]byte(password), hash3) != nil {
		t.Fatal("password verification failed")
		t.FailNow()
	}
}

func TestHashVerifyBase64(t *testing.T) {
	const password = "my super secret password"
	hash0, err := passhashx.HashBase64([]byte(password), passhashx.SecurityLevelHigh)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	if passhashx.VerifyBase64([]byte(password), hash0) != nil {
		t.Fatal("password verification failed")
		t.FailNow()
	}

	hash1, err := passhashx.HashBase64([]byte(password), passhashx.SecurityLevelLow)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	if passhashx.VerifyBase64([]byte(password), hash1) != nil {
		t.Fatal("password verification failed")
		t.FailNow()
	}

	hash2, err := passhashx.HashBase64([]byte(password), passhashx.SecurityLevelMobileHigh)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	if passhashx.VerifyBase64([]byte(password), hash2) != nil {
		t.Fatal("password verification failed")
		t.FailNow()
	}
}
