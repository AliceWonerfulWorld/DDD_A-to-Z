package config

import "testing"

func TestEnvOrDefault(t *testing.T) {
	t.Run("環境変数が空ならデフォルト値を返す", func(t *testing.T) {
		t.Setenv("LANG_WAR_TEST_VALUE", "")

		if got := EnvOrDefault("LANG_WAR_TEST_VALUE", "fallback"); got != "fallback" {
			t.Fatalf("EnvOrDefault の戻り値 = %q, 期待値 fallback", got)
		}
	})

	t.Run("環境変数が設定されていればその値を返す", func(t *testing.T) {
		t.Setenv("LANG_WAR_TEST_VALUE", "configured")

		if got := EnvOrDefault("LANG_WAR_TEST_VALUE", "fallback"); got != "configured" {
			t.Fatalf("EnvOrDefault の戻り値 = %q, 期待値 configured", got)
		}
	})
}
