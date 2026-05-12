package config

import "testing"

func TestEnvOrDefault(t *testing.T) {
	t.Setenv("LANG_WAR_TEST_VALUE", "")

	if got := EnvOrDefault("LANG_WAR_TEST_VALUE", "fallback"); got != "fallback" {
		t.Fatalf("EnvOrDefault returned %q, want fallback", got)
	}

	t.Setenv("LANG_WAR_TEST_VALUE", "configured")

	if got := EnvOrDefault("LANG_WAR_TEST_VALUE", "fallback"); got != "configured" {
		t.Fatalf("EnvOrDefault returned %q, want configured", got)
	}
}
