package ctrx

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
)

func EncryptCTR(data []byte, key []byte) ([]byte, error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}
	iv := make([]byte, block.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCTR(block, iv)
	result := make([]byte, len(data))
	mode.XORKeyStream(result, data)
	return append(iv, result...), nil
}

func DecryptCTR(data []byte, key []byte) ([]byte, error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}
	ivSize := block.BlockSize()
	iv, ciphertext := data[:ivSize], data[ivSize:]
	mode := cipher.NewCTR(block, iv)
	result := make([]byte, len(ciphertext))
	mode.XORKeyStream(result, ciphertext)
	return result, nil
}
