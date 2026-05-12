package config

import (
	"fmt"
	"os"
)

const defaultGitHubRedirectURL = "http://localhost:8080/auth/github/callback"

type GitHubOAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func GitHubOAuthFromEnv() (GitHubOAuth, error) {
	config := GitHubOAuth{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
	}
	if config.ClientID == "" || config.ClientSecret == "" || config.RedirectURL == "" || config.RedirectURL == defaultGitHubRedirectURL {
		return GitHubOAuth{}, fmt.Errorf("GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET, and non-localhost GITHUB_REDIRECT_URL must be set")
	}

	return config, nil
}
