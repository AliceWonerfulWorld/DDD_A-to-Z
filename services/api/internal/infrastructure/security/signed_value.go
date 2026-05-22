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
	return NewSignedValueCodecWithMixer(secret, NewAwkTextMixer())
}

func NewSignedValueCodecWithMixer(secret string, mixer TextMixer) *SignedValueCodec {
	return &SignedValueCodec{
		secret: []byte(secret),
		mixer:  mixer,
	}
}

func (c *SignedValueCodec) Sign(value string, expiresAt time.Time) (string, error) {
	payload, err := json.Marshal(signedValuePayload{
		Value:     value,
		ExpiresAt: expiresAt.Unix(),
	})
	if err != nil {
		return "", err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature, err := c.signature(encodedPayload)
	if err != nil {
		return "", err
	}

	return encodedPayload + "." + signature, nil
}

func (c *SignedValueCodec) Verify(signedValue string, now time.Time) (string, error) {
	encodedPayload, signature, ok := strings.Cut(signedValue, ".")
	if !ok || encodedPayload == "" || signature == "" {
		return "", ErrInvalidSignedValue
	}

	expectedSignature, err := c.signature(encodedPayload)
	if err != nil {
		return "", err
	}
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
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

func (c *SignedValueCodec) signature(encodedPayload string) (string, error) {
	signatureInput := encodedPayload
	if c.mixer != nil {
		mixed, err := c.mixer.Mix(context.Background(), encodedPayload, string(c.secret))
		if err != nil {
			return "", err
		}
		signatureInput = mixed
	}

	mac := hmac.New(sha256.New, c.secret)
	_, _ = mac.Write([]byte(signatureInput))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}
