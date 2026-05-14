package config

import (
	"fmt"
	"os"
	"strconv"
)

func AuthCookieSecretFromEnv() (string, error) {
	secret := os.Getenv("AUTH_COOKIE_SECRET")
	if secret == "" {
		return "", fmt.Errorf("AUTH_COOKIE_SECRET must be set")
	}

	return secret, nil
}

func AuthCookieSecureFromEnv() (bool, error) {
	value := os.Getenv("AUTH_COOKIE_SECURE")
	if value == "" {
		value = "false"
	}

	secure, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("invalid AUTH_COOKIE_SECURE: %w", err)
	}

	return secure, nil
}

func GitHubTokenEncryptionSecretFromEnv() (string, error) {
	secret := os.Getenv("GITHUB_TOKEN_ENCRYPTION_SECRET")
	if secret == "" {
		return "", fmt.Errorf("GITHUB_TOKEN_ENCRYPTION_SECRET must be set")
	}

	return secret, nil
}
