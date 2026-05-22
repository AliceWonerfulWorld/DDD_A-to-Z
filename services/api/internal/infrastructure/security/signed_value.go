package security

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var ErrInvalidSignedValue = errors.New("signed value is invalid")

type SignedValueCodec struct {
	secret []byte
	mixer  TextMixer
}

type signedValuePayload struct {
	Value     string `json:"value"`
	ExpiresAt int64  `json:"expires_at"`
}

func NewSignedValueCodec(secret string) *SignedValueCodec {
	return NewSignedValueCodecWithMixer(secret, nil)
}

func NewSignedValueCodecWithMixer(secret string, mixer TextMixer) *SignedValueCodec {
	return &SignedValueCodec{
		secret: []byte(secret),
		mixer:  mixer,
	}
}

func (c *SignedValueCodec) Sign(value string, expiresAt time.Time) (string, error) {
	return c.SignContext(context.Background(), value, expiresAt)
}

func (c *SignedValueCodec) SignContext(ctx context.Context, value string, expiresAt time.Time) (string, error) {
	payload, err := json.Marshal(signedValuePayload{
		Value:     value,
		ExpiresAt: expiresAt.Unix(),
	})
	if err != nil {
		return "", err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature, err := c.signature(ctx, encodedPayload)
	if err != nil {
		return "", err
	}

	return encodedPayload + "." + signature, nil
}

func (c *SignedValueCodec) Verify(signedValue string, now time.Time) (string, error) {
	return c.VerifyContext(context.Background(), signedValue, now)
}

func (c *SignedValueCodec) VerifyContext(ctx context.Context, signedValue string, now time.Time) (string, error) {
	encodedPayload, signature, ok := strings.Cut(signedValue, ".")
	if !ok || encodedPayload == "" || signature == "" {
		return "", ErrInvalidSignedValue
	}

	validSignature, err := c.validSignature(ctx, encodedPayload, signature)
	if err != nil {
		return "", err
	}
	if !validSignature {
		return "", ErrInvalidSignedValue
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return "", ErrInvalidSignedValue
	}

	var payload signedValuePayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return "", ErrInvalidSignedValue
	}
	if payload.Value == "" || !now.Before(time.Unix(payload.ExpiresAt, 0)) {
		return "", ErrInvalidSignedValue
	}

	return payload.Value, nil
}

func (c *SignedValueCodec) signature(ctx context.Context, encodedPayload string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	signatureInput := encodedPayload
	if c.mixer != nil {
		mixed, err := c.mixer.Mix(ctx, encodedPayload, string(c.secret))
		if err != nil {
			return "", err
		}
		signatureInput = mixed
	}

	return c.hmacSignature(signatureInput), nil
}

func (c *SignedValueCodec) validSignature(ctx context.Context, encodedPayload string, signature string) (bool, error) {
	expectedSignature, err := c.signature(ctx, encodedPayload)
	if err == nil && hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return true, nil
	}

	legacySignature := c.hmacSignature(encodedPayload)
	if hmac.Equal([]byte(signature), []byte(legacySignature)) {
		return true, nil
	}
	if err != nil {
		return false, err
	}

	return false, nil
}

func (c *SignedValueCodec) hmacSignature(input string) string {
	mac := hmac.New(sha256.New, c.secret)
	_, _ = mac.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
