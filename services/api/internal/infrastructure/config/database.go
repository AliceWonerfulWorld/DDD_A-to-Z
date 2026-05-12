package config

func DatabaseURLFromEnv() string {
	return EnvOrDefault("DATABASE_URL", "postgres://lang_war:lang_war_password@localhost:5432/lang_war?sslmode=disable")
}
