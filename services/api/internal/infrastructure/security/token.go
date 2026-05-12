package security

import (
	"crypto/rand"
	"encoding/base64"
)

type SecureTokenGenerator struct{}

func NewSecureTokenGenerator() *SecureTokenGenerator {
	return &SecureTokenGenerator{}
}

func (g *SecureTokenGenerator) NewToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
