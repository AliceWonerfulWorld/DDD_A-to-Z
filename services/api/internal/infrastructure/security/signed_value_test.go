package security

import (
	"testing"
	"time"
)

func TestSignedValueCodecVerify(t *testing.T) {
	codec := NewSignedValueCodec("test-secret")
	expiresAt := time.Date(2026, 5, 12, 12, 10, 0, 0, time.UTC)

	signedValue, err := codec.Sign("state-token", expiresAt)
	if err != nil {
		t.Fatalf("Sign returned error: %v", err)
	}

	got, err := codec.Verify(signedValue, time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if got != "state-token" {
		t.Fatalf("Verify returned %q, want state-token", got)
	}
}

func TestSignedValueCodecRejectsExpiredValue(t *testing.T) {
	codec := NewSignedValueCodec("test-secret")
	expiresAt := time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC)

	signedValue, err := codec.Sign("state-token", expiresAt)
	if err != nil {
		t.Fatalf("Sign returned error: %v", err)
	}

	if _, err := codec.Verify(signedValue, time.Date(2026, 5, 12, 12, 1, 0, 0, time.UTC)); err != ErrInvalidSignedValue {
		t.Fatalf("Verify error = %v, want ErrInvalidSignedValue", err)
	}
}
