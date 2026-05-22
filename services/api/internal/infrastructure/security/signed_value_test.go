package security

import (
	"context"
	"errors"
	"os/exec"
	"testing"
	"time"
)

type failingTextMixer struct{}

func (failingTextMixer) Mix(context.Context, string, string) (string, error) {
	return "", errors.New("mixer failed")
}

type contextCheckingTextMixer struct {
	want context.Context
}

type signedValueTestContextKey struct{}

func (m contextCheckingTextMixer) Mix(ctx context.Context, input string, _ string) (string, error) {
	if ctx != m.want {
		return "", errors.New("unexpected context")
	}

	return input, nil
}

func TestSignedValueCodecVerify(t *testing.T) {
	requireAwkCommand(t)

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
	requireAwkCommand(t)

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

func TestSignedValueCodecPassesContextToMixer(t *testing.T) {
	ctx := context.WithValue(context.Background(), signedValueTestContextKey{}, "request-context")
	codec := NewSignedValueCodecWithMixer("test-secret", contextCheckingTextMixer{want: ctx})
	expiresAt := time.Date(2026, 5, 12, 12, 10, 0, 0, time.UTC)

	signedValue, err := codec.SignContext(ctx, "state-token", expiresAt)
	if err != nil {
		t.Fatalf("SignContext がエラーを返しました: %v", err)
	}

	if _, err := codec.VerifyContext(ctx, signedValue, time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC)); err != nil {
		t.Fatalf("VerifyContext がエラーを返しました: %v", err)
	}
}

func TestSignedValueCodecVerifiesLegacySignature(t *testing.T) {
	requireAwkCommand(t)

	legacyCodec := NewSignedValueCodecWithMixer("test-secret", nil)
	codec := NewSignedValueCodec("test-secret")
	expiresAt := time.Date(2026, 5, 12, 12, 10, 0, 0, time.UTC)

	signedValue, err := legacyCodec.Sign("state-token", expiresAt)
	if err != nil {
		t.Fatalf("legacy Sign がエラーを返しました: %v", err)
	}

	got, err := codec.Verify(signedValue, time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Verify が legacy signature でエラーを返しました: %v", err)
	}
	if got != "state-token" {
		t.Fatalf("Verify の戻り値 = %q, 期待値 state-token", got)
	}
}

func TestSignedValueCodecVerifiesLegacySignatureWhenMixerFails(t *testing.T) {
	legacyCodec := NewSignedValueCodecWithMixer("test-secret", nil)
	codec := NewSignedValueCodecWithMixer("test-secret", failingTextMixer{})
	expiresAt := time.Date(2026, 5, 12, 12, 10, 0, 0, time.UTC)

	signedValue, err := legacyCodec.Sign("state-token", expiresAt)
	if err != nil {
		t.Fatalf("legacy Sign がエラーを返しました: %v", err)
	}

	got, err := codec.Verify(signedValue, time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Verify が legacy signature でエラーを返しました: %v", err)
	}
	if got != "state-token" {
		t.Fatalf("Verify の戻り値 = %q, 期待値 state-token", got)
	}
}

func TestSignedValueCodecReturnsMixerError(t *testing.T) {
	codec := NewSignedValueCodecWithMixer("test-secret", failingTextMixer{})
	expiresAt := time.Date(2026, 5, 12, 12, 10, 0, 0, time.UTC)

	if _, err := codec.Sign("state-token", expiresAt); err == nil {
		t.Fatal("Sign のエラー = nil, 期待値 mixer error")
	}
}

func requireAwkCommand(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("awk"); err != nil {
		t.Skip("awk command is not available")
	}
}
