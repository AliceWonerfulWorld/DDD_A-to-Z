package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	githubapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/github"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/repositoryanalysis"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

const (
	gitHubAPIBaseURL = "https://api.github.com"
	gitHubAPIVersion = "2022-11-28"
	userAgent        = "lang-war-api"
)

type RepositoryClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewRepositoryClient(httpClient *http.Client) *RepositoryClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &RepositoryClient{
		httpClient: httpClient,
		baseURL:    gitHubAPIBaseURL,
	}
}

func (c *RepositoryClient) ListRepositories(ctx context.Context, accessToken string, userID user.ID, syncedAt time.Time) ([]repositoryanalysis.Repository, error) {
	nextURL := c.baseURL + "/user/repos?affiliation=owner,collaborator,organization_member&per_page=100&sort=updated&direction=desc"
	repositories := []repositoryanalysis.Repository{}

	for nextURL != "" {
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, nextURL, nil)
		if err != nil {
			return nil, err
		}
		request.Header.Set("Accept", "application/vnd.github+json")
		request.Header.Set("Authorization", "Bearer "+accessToken)
		request.Header.Set("User-Agent", userAgent)
		request.Header.Set("X-GitHub-Api-Version", gitHubAPIVersion)

		response, err := c.httpClient.Do(request)
		if err != nil {
			return nil, classifyTransportError(err)
		}
		payload, closeErr := decodeRepositoryPage(response)
		if closeErr != nil {
			return nil, closeErr
		}

		for _, item := range payload {
			repository, err := item.toDomain(userID, syncedAt)
			if err != nil {
				return nil, err
			}
			repositories = append(repositories, repository)
		}

		nextURL = nextLink(response.Header.Get("Link"))
	}

	return repositories, nil
}

func decodeRepositoryPage(response *http.Response) ([]repositoryPayload, error) {
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, classifyAPIError(response)
	}

	var payload []repositoryPayload
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload, nil
}

type repositoryPayload struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	FullName      string     `json:"full_name"`
	Private       bool       `json:"private"`
	Fork          bool       `json:"fork"`
	Archived      bool       `json:"archived"`
	DefaultBranch string     `json:"default_branch"`
	Language      *string    `json:"language"`
	HTMLURL       string     `json:"html_url"`
	PushedAt      *time.Time `json:"pushed_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Owner         struct {
		Login string `json:"login"`
	} `json:"owner"`
}

func (p repositoryPayload) toDomain(userID user.ID, syncedAt time.Time) (repositoryanalysis.Repository, error) {
	language := ""
	if p.Language != nil {
		language = *p.Language
	}

	return repositoryanalysis.NewRepository(repositoryanalysis.Repository{
		GitHubID:        p.ID,
		UserID:          userID,
		Owner:           p.Owner.Login,
		Name:            p.Name,
		FullName:        p.FullName,
		Private:         p.Private,
		Fork:            p.Fork,
		Archived:        p.Archived,
		DefaultBranch:   p.DefaultBranch,
		Language:        language,
		HTMLURL:         p.HTMLURL,
		PushedAt:        p.PushedAt,
		GitHubUpdatedAt: p.UpdatedAt,
		SyncedAt:        syncedAt,
	})
}

type errorPayload struct {
	Message string `json:"message"`
}

func classifyAPIError(response *http.Response) error {
	payload := errorPayload{}
	_ = json.NewDecoder(response.Body).Decode(&payload)
	message := payload.Message
	if message == "" {
		message = fmt.Sprintf("github api request failed: status %d", response.StatusCode)
	}

	if isRateLimited(response, message) {
		retryAfter, resetAt := retryWindow(response)
		if retryAfter == 0 && resetAt == nil {
			retryAfter = time.Minute
		}
		return &githubapp.APIError{
			Kind:       githubapp.ErrorKindRateLimited,
			Message:    "GitHub API rate limit exceeded",
			StatusCode: response.StatusCode,
			RetryAfter: retryAfter,
			ResetAt:    resetAt,
		}
	}

	switch response.StatusCode {
	case http.StatusUnauthorized:
		return &githubapp.APIError{Kind: githubapp.ErrorKindTokenInvalid, Message: message, StatusCode: response.StatusCode}
	case http.StatusForbidden:
		return &githubapp.APIError{Kind: githubapp.ErrorKindPermissionDenied, Message: message, StatusCode: response.StatusCode}
	case http.StatusNotFound:
		return &githubapp.APIError{Kind: githubapp.ErrorKindPermissionDeniedOrNotFound, Message: message, StatusCode: response.StatusCode}
	}
	if response.StatusCode >= http.StatusInternalServerError {
		return &githubapp.APIError{Kind: githubapp.ErrorKindUnavailable, Message: message, StatusCode: response.StatusCode}
	}

	return fmt.Errorf("github api request failed: status %d: %s", response.StatusCode, message)
}

func classifyTransportError(err error) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return &githubapp.APIError{Kind: githubapp.ErrorKindUnavailable, Message: "GitHub API unavailable", StatusCode: http.StatusBadGateway}
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return &githubapp.APIError{Kind: githubapp.ErrorKindUnavailable, Message: "GitHub API unavailable", StatusCode: http.StatusBadGateway}
	}

	return err
}

func isRateLimited(response *http.Response, message string) bool {
	if response.StatusCode == http.StatusTooManyRequests {
		return true
	}

	statusCanRateLimit := response.StatusCode == http.StatusForbidden || response.StatusCode == http.StatusTooManyRequests
	if !statusCanRateLimit {
		return false
	}
	if response.Header.Get("X-RateLimit-Remaining") == "0" {
		return true
	}
	if response.Header.Get("Retry-After") != "" {
		return true
	}

	return strings.Contains(strings.ToLower(message), "secondary rate limit")
}

func retryWindow(response *http.Response) (time.Duration, *time.Time) {
	if retryAfterHeader := response.Header.Get("Retry-After"); retryAfterHeader != "" {
		seconds, err := strconv.ParseInt(retryAfterHeader, 10, 64)
		if err == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second, nil
		}
	}

	resetHeader := response.Header.Get("X-RateLimit-Reset")
	if resetHeader == "" {
		return 0, nil
	}
	epoch, err := strconv.ParseInt(resetHeader, 10, 64)
	if err != nil || epoch <= 0 {
		return 0, nil
	}
	resetAt := time.Unix(epoch, 0).UTC()

	return 0, &resetAt
}

func nextLink(linkHeader string) string {
	for _, link := range strings.Split(linkHeader, ",") {
		section := strings.TrimSpace(link)
		if !strings.Contains(section, `rel="next"`) {
			continue
		}

		start := strings.Index(section, "<")
		end := strings.Index(section, ">")
		if start == -1 || end == -1 || end <= start+1 {
			return ""
		}

		nextURL := section[start+1 : end]
		if _, err := url.ParseRequestURI(nextURL); err != nil {
			return ""
		}
		return nextURL
	}

	return ""
}
