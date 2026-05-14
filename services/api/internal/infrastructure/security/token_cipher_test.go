package security

import (
	"errors"
	"testing"
)

func TestTokenCipher(t *testing.T) {
	t.Run("同じ associated data で復号できる", func(t *testing.T) {
		cipher, err := NewTokenCipher("test-secret")
		if err != nil {
			t.Fatalf("NewTokenCipher() がエラーを返しました: %v", err)
		}

		ciphertext, err := cipher.Encrypt("github-token", "user_1")
		if err != nil {
			t.Fatalf("Encrypt() がエラーを返しました: %v", err)
		}

		plaintext, err := cipher.Decrypt(ciphertext, "user_1")
		if err != nil {
			t.Fatalf("Decrypt() がエラーを返しました: %v", err)
		}
		if plaintext != "github-token" {
			t.Fatalf("plaintext = %q, 期待値 github-token", plaintext)
		}
	})

	t.Run("違う associated data では復号できない", func(t *testing.T) {
		cipher, err := NewTokenCipher("test-secret")
		if err != nil {
			t.Fatalf("NewTokenCipher() がエラーを返しました: %v", err)
		}

		ciphertext, err := cipher.Encrypt("github-token", "user_1")
		if err != nil {
			t.Fatalf("Encrypt() がエラーを返しました: %v", err)
		}

		if _, err := cipher.Decrypt(ciphertext, "user_2"); !errors.Is(err, ErrInvalidCiphertext) {
			t.Fatalf("Decrypt() error = %v, 期待値 ErrInvalidCiphertext", err)
		}
	})
}
