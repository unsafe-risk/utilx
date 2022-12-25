package cbcx

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
)

func Encrypt(data []byte, key []byte) ([]byte, error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}
	iv := make([]byte, block.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	result := make([]byte, len(data))
	mode.CryptBlocks(result, data)
	return append(iv, result...), nil
}

func Decrypt(data []byte, key []byte) ([]byte, error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}
	ivSize := block.BlockSize()
	iv, ciphertext := data[:ivSize], data[ivSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	result := make([]byte, len(ciphertext))
	mode.CryptBlocks(result, ciphertext)
	return result, nil
}
