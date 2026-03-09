package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}

func ConvertToAESKeyFromRMK(rmk string) []byte {
	hash := sha256.Sum256([]byte(rmk))

	return hash[:]
}

func GenerateEncryptedKEK(rmk string) ([]byte, error) {
	kek, err := GenerateKey() // random 32-byte KEK
	if err != nil {
		return nil, err
	}

	rmkKey := ConvertToAESKeyFromRMK(rmk)

	return EncryptKEK(kek, rmkKey)
}

func EncryptKEK(kek, rmk []byte) ([]byte, error) {
	// Cipher Block
	cipherBlock, err := aes.NewCipher(rmk)
	if err != nil {
		return nil, err
	}

	// GCM Block

	gcm, err := cipher.NewGCM(cipherBlock)

	if err != nil {
		return nil, err
	}

	// Generate Nonce
	nonce := make([]byte, gcm.NonceSize())

	// Cipher Text
	cipherText := gcm.Seal(nonce, nonce, kek, nil)

	return cipherText, nil
}

func DecryptKEK(encryptedKEK, rmk []byte) ([]byte, error) {
	cipherBlock, err := aes.NewCipher(rmk)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(cipherBlock)

	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	if len(encryptedKEK) < nonceSize {
		return nil, fmt.Errorf("invalid encrypted kek")
	}

	nonce := encryptedKEK[:nonceSize]
	cipherText := encryptedKEK[nonceSize:]

	return gcm.Open(nil, nonce, cipherText, nil)
}

func Encrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func Decrypt(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	if len(cipherText) < nonceSize {
		return nil, fmt.Errorf("invalid cipherText")
	}

	nonce := cipherText[:nonceSize]
	encrypted := cipherText[nonceSize:]

	return gcm.Open(nil, nonce, encrypted, nil)
}

// -- Generate Password
func GeneratePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"

	b := make([]byte, length)

	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))

		if err != nil {
			return "", err
		}

		b[i] = charset[n.Int64()]
	}

	return string(b), nil
}
