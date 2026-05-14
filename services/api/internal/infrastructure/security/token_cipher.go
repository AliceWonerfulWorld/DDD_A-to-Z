package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var ErrInvalidCiphertext = errors.New("ciphertext is invalid")

type TokenCipher struct {
	aead cipher.AEAD
}

func NewTokenCipher(secret string) (*TokenCipher, error) {
	if secret == "" {
		return nil, errors.New("token encryption secret is required")
	}

	key := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &TokenCipher{aead: aead}, nil
}

func (c *TokenCipher) Encrypt(plaintext, associatedData string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	nonce := make([]byte, c.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	sealed := c.aead.Seal(nonce, nonce, []byte(plaintext), []byte(associatedData))
	return base64.RawURLEncoding.EncodeToString(sealed), nil
}

func (c *TokenCipher) Decrypt(ciphertext, associatedData string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	sealed, err := base64.RawURLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", ErrInvalidCiphertext
	}
	if len(sealed) < c.aead.NonceSize() {
		return "", ErrInvalidCiphertext
	}

	nonce := sealed[:c.aead.NonceSize()]
	payload := sealed[c.aead.NonceSize():]
	plaintext, err := c.aead.Open(nil, nonce, payload, []byte(associatedData))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidCiphertext, err)
	}

	return string(plaintext), nil
}
