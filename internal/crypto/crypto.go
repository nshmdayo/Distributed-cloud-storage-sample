// Package crypto provides cryptographic functions for the distributed storage system
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// EncryptionKey represents an encryption key
type EncryptionKey []byte

// GenerateKey generates a new AES-256 encryption key
func GenerateKey() (EncryptionKey, error) {
	key := make([]byte, 32) // 256 bits
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return EncryptionKey(key), nil
}

// DeriveKey derives an encryption key from a password using SHA-256
func DeriveKey(password string) EncryptionKey {
	hash := sha256.Sum256([]byte(password))
	return EncryptionKey(hash[:])
}

// Encrypt encrypts data using AES-256-GCM
func Encrypt(data []byte, key EncryptionKey) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt decrypts data using AES-256-GCM
func Decrypt(ciphertext []byte, key EncryptionKey) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Hash calculates SHA-256 hash of data
func Hash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// VerifyHash verifies if data matches the expected hash
func VerifyHash(data []byte, expectedHash []byte) bool {
	actualHash := Hash(data)
	if len(actualHash) != len(expectedHash) {
		return false
	}

	for i := range actualHash {
		if actualHash[i] != expectedHash[i] {
			return false
		}
	}
	return true
}
