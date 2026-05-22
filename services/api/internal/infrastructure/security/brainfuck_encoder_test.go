package security

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
)

func TestBrainfuckEncoder(t *testing.T) {
	encoder := NewBrainfuckEncoder()
	input := "awk-mixed-token::15:4"

	encoded, err := encoder.Encode(input)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if encoded == input {
		t.Fatal("Encode() returned the original input")
	}
	if _, err := base64.RawURLEncoding.DecodeString(encoded); err != nil {
		t.Fatalf("Encode() returned non base64url payload: %v", err)
	}

	decoded, err := encoder.Decode(encoded)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if decoded != input {
		t.Fatalf("Decode() = %q, want %q", decoded, input)
	}
	if encoder.Type() != EncoderTypeBrainfuck {
		t.Fatalf("Type() = %q, want %q", encoder.Type(), EncoderTypeBrainfuck)
	}
}

func TestBrainfuckEncoderIsDeterministic(t *testing.T) {
	encoder := NewBrainfuckEncoder()

	first, err := encoder.Encode("state-token")
	if err != nil {
		t.Fatalf("first Encode() error = %v", err)
	}
	second, err := encoder.Encode("state-token")
	if err != nil {
		t.Fatalf("second Encode() error = %v", err)
	}

	if first != second {
		t.Fatalf("Encode() is not deterministic: first %q, second %q", first, second)
	}
}

func TestBrainfuckEncoderDecodeInvalidPayload(t *testing.T) {
	encoder := NewBrainfuckEncoder()

	if _, err := encoder.Decode("not+base64url"); err == nil {
		t.Fatal("Decode() error = nil, want error")
	}
}

type staticTokenSource struct {
	token string
	err   error
}

func (s staticTokenSource) NewToken() (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.token, nil
}

type fakeTextMixer struct {
	input string
	salt  string
	err   error
}

func (m *fakeTextMixer) Mix(ctx context.Context, input string, salt string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	m.input = input
	m.salt = salt
	if m.err != nil {
		return "", m.err
	}
	return "mixed-" + input + "-" + salt, nil
}

func TestObfuscatedTokenGenerator(t *testing.T) {
	mixer := &fakeTextMixer{}
	generator := NewObfuscatedTokenGenerator(staticTokenSource{token: "raw-token"}, mixer, "salt")

	token, err := generator.NewToken()
	if err != nil {
		t.Fatalf("NewToken() error = %v", err)
	}
	if mixer.input != "raw-token" {
		t.Fatalf("mixer input = %q, want raw-token", mixer.input)
	}
	if mixer.salt != "salt" {
		t.Fatalf("mixer salt = %q, want salt", mixer.salt)
	}

	decoded, err := NewBrainfuckEncoder().Decode(token)
	if err != nil {
		t.Fatalf("Decode(obfuscated token) error = %v", err)
	}
	if decoded != "mixed-raw-token-salt" {
		t.Fatalf("decoded token = %q, want mixed-raw-token-salt", decoded)
	}
}

func TestObfuscatedTokenGeneratorReturnsSourceError(t *testing.T) {
	want := errors.New("source failed")
	generator := NewObfuscatedTokenGenerator(staticTokenSource{err: want}, &fakeTextMixer{}, "salt")

	if _, err := generator.NewToken(); !errors.Is(err, want) {
		t.Fatalf("NewToken() error = %v, want %v", err, want)
	}
}

func TestObfuscatedTokenGeneratorReturnsMixerError(t *testing.T) {
	want := errors.New("mixer failed")
	generator := NewObfuscatedTokenGenerator(staticTokenSource{token: "raw"}, &fakeTextMixer{err: want}, "salt")

	if _, err := generator.NewToken(); !errors.Is(err, want) {
		t.Fatalf("NewToken() error = %v, want %v", err, want)
	}
}
