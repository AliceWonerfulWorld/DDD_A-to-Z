package security

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/brainfuck"
)

type EncoderType string

const EncoderTypeBrainfuck EncoderType = "brainfuck"

const (
	brainfuckShiftEncode = ",>+++++++++++++[<+>-]<."
	brainfuckShiftDecode = ",>+++++++++++++[<->-]<."
)

type BrainfuckEncoder struct {
	encoderType EncoderType
}

func NewBrainfuckEncoder() *BrainfuckEncoder {
	return &BrainfuckEncoder{encoderType: EncoderTypeBrainfuck}
}

func (e *BrainfuckEncoder) Type() EncoderType {
	return e.encoderType
}

func (e *BrainfuckEncoder) Encode(input string) (string, error) {
	encoded, err := e.transform(input, brainfuckShiftEncode)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString([]byte(encoded)), nil
}

func (e *BrainfuckEncoder) Decode(input string) (string, error) {
	decodedBytes, err := base64.RawURLEncoding.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("decode brainfuck payload: %w", err)
	}

	return e.transform(string(decodedBytes), brainfuckShiftDecode)
}

func (e *BrainfuckEncoder) transform(input string, code string) (string, error) {
	output := make([]byte, 0, len(input))
	for _, b := range []byte(input) {
		transformed, err := brainfuck.ExecuteBF(code, string([]byte{b}))
		if err != nil {
			return "", fmt.Errorf("%s encode byte: %w", e.encoderType, err)
		}
		if len(transformed) != 1 {
			return "", fmt.Errorf("%s encode byte: expected 1 output byte, got %d", e.encoderType, len(transformed))
		}
		output = append(output, transformed[0])
	}

	return string(output), nil
}

type ObfuscatedTokenGenerator struct {
	source  interface{ NewToken() (string, error) }
	mixer   TextMixer
	encoder *BrainfuckEncoder
	salt    string
}

func NewObfuscatedTokenGenerator(source interface{ NewToken() (string, error) }, mixer TextMixer, salt string) *ObfuscatedTokenGenerator {
	return &ObfuscatedTokenGenerator{
		source:  source,
		mixer:   mixer,
		encoder: NewBrainfuckEncoder(),
		salt:    salt,
	}
}

func (g *ObfuscatedTokenGenerator) NewToken() (string, error) {
	token, err := g.source.NewToken()
	if err != nil {
		return "", err
	}

	mixed, err := g.mixer.Mix(context.Background(), token, g.salt)
	if err != nil {
		return "", err
	}

	return g.encoder.Encode(mixed)
}
