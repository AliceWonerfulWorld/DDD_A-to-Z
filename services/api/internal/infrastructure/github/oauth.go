package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/infrastructure/config"
	"golang.org/x/oauth2"
	oauthgithub "golang.org/x/oauth2/github"
)

type OAuthClient struct {
	config     config.GitHubOAuth
	httpClient *http.Client
}

func NewOAuthClient(config config.GitHubOAuth, httpClient *http.Client) *OAuthClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &OAuthClient{
		config:     config,
		httpClient: httpClient,
	}
}

func (c *OAuthClient) AuthCodeURL(state string) string {
	return c.oauthConfig().AuthCodeURL(state)
}

func (c *OAuthClient) ExchangeProfile(ctx context.Context, code string) (user.GitHubProfile, error) {
	token, err := c.oauthConfig().Exchange(c.contextWithHTTPClient(ctx), code)
	if err != nil {
		return user.GitHubProfile{}, err
	}

	return c.fetchProfile(ctx, token.AccessToken)
}

func (c *OAuthClient) oauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		RedirectURL:  c.config.RedirectURL,
		Scopes:       []string{"read:user"},
		Endpoint:     oauthgithub.Endpoint,
	}
}

func (c *OAuthClient) contextWithHTTPClient(ctx context.Context) context.Context {
	return context.WithValue(ctx, oauth2.HTTPClient, c.httpClient)
}

func (c *OAuthClient) fetchProfile(ctx context.Context, accessToken string) (user.GitHubProfile, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return user.GitHubProfile{}, err
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return user.GitHubProfile{}, err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return user.GitHubProfile{}, fmt.Errorf("github user request failed: status %d", response.StatusCode)
	}

	var payload struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return user.GitHubProfile{}, err
	}

	return user.GitHubProfile{
		GitHubID:  payload.ID,
		Username:  payload.Login,
		AvatarURL: payload.AvatarURL,
	}, nil
}
