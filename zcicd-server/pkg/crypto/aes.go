package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// Encryptor provides AES-256-GCM encryption/decryption.
// Format: version(1B) + nonce(12B) + ciphertext
type Encryptor struct {
	key []byte
}

const (
	version   = byte(0x01)
	nonceSize = 12
)

func NewEncryptor(key string) (*Encryptor, error) {
	k := []byte(key)
	if len(k) != 32 {
		return nil, fmt.Errorf("AES key must be 32 bytes, got %d", len(k))
	}
	return &Encryptor{key: k}, nil
}

func (e *Encryptor) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	result := make([]byte, 0, 1+nonceSize+len(ciphertext))
	result = append(result, version)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

func (e *Encryptor) Decrypt(data []byte) ([]byte, error) {
	if len(data) < 1+nonceSize+16 { // 16 = GCM tag size
		return nil, errors.New("ciphertext too short")
	}

	if data[0] != version {
		return nil, fmt.Errorf("unsupported encryption version: %d", data[0])
	}

	nonce := data[1 : 1+nonceSize]
	ciphertext := data[1+nonceSize:]

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, ciphertext, nil)
}
