package config

func AuthCookieSecretFromEnv() string {
	return EnvOrDefault("AUTH_COOKIE_SECRET", "development-only-auth-cookie-secret")
}

func AuthCookieSecureFromEnv() bool {
	return EnvOrDefault("AUTH_COOKIE_SECURE", "false") == "true"
}
