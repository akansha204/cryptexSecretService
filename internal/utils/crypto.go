package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	aesgcm    cipher.AEAD
	nonceSize int
)

func init() {
	if err := initCipherFromEnv(); err != nil {
		panic(fmt.Sprintf("cryptoutil initialization failed: %v", err))
	}
}

func initCipherFromEnv() error {
	raw := os.Getenv("SECRET_ENCRYPTION_KEY")
	if raw == "" {
		return errors.New("SECRET_ENCRYPTION_KEY is not set")
	}

	key, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		key = []byte(raw)
	}

	if len(key) != 32 {
		return fmt.Errorf("invalid encryption key length: got %d bytes, need 32 bytes (256 bits). If using base64, ensure it decodes to 32 bytes", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}

	aesgcm = gcm
	nonceSize = gcm.NonceSize()
	return nil
}

func Encrypt(plaintext string) (string, error) {
	if aesgcm == nil {
		return "", errors.New("encryption not initialized")
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("cannot generate nonce: %w", err)
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

	out := append(nonce, ciphertext...)

	encoded := base64.StdEncoding.EncodeToString(out)
	return encoded, nil
}

func Decrypt(ciphertextB64 string) (string, error) {
	if aesgcm == nil {
		return "", errors.New("decryption not initialized")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", fmt.Errorf("invalid base64 ciphertext: %w", err)
	}

	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed or data tampered: %w", err)
	}

	return string(plaintext), nil
}
