package security

import (
	"testing"
	"time"
)

func TestSignedValueCodecVerify(t *testing.T) {
	t.Run("署名付き値を検証できる", func(t *testing.T) {
		codec := NewSignedValueCodec("test-secret")
		expiresAt := time.Date(2026, 5, 12, 12, 10, 0, 0, time.UTC)

		signedValue, err := codec.Sign("state-token", expiresAt)
		if err != nil {
			t.Fatalf("Sign がエラーを返しました: %v", err)
		}

		got, err := codec.Verify(signedValue, time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC))
		if err != nil {
			t.Fatalf("Verify がエラーを返しました: %v", err)
		}
		if got != "state-token" {
			t.Fatalf("Verify の戻り値 = %q, 期待値 state-token", got)
		}
	})
}

func TestSignedValueCodecRejectsExpiredValue(t *testing.T) {
	t.Run("期限切れの署名付き値を拒否する", func(t *testing.T) {
		codec := NewSignedValueCodec("test-secret")
		expiresAt := time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC)

		signedValue, err := codec.Sign("state-token", expiresAt)
		if err != nil {
			t.Fatalf("Sign がエラーを返しました: %v", err)
		}

		if _, err := codec.Verify(signedValue, time.Date(2026, 5, 12, 12, 1, 0, 0, time.UTC)); err != ErrInvalidSignedValue {
			t.Fatalf("Verify のエラー = %v, 期待値 ErrInvalidSignedValue", err)
		}
	})
}
