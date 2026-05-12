package config

type GitHubOAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func GitHubOAuthFromEnv() GitHubOAuth {
	return GitHubOAuth{
		ClientID:     EnvOrDefault("GITHUB_CLIENT_ID", ""),
		ClientSecret: EnvOrDefault("GITHUB_CLIENT_SECRET", ""),
		RedirectURL:  EnvOrDefault("GITHUB_REDIRECT_URL", "http://localhost:8080/auth/github/callback"),
	}
}
