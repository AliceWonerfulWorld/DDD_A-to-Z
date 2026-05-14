// Package github owns GitHub repository synchronization use cases.
package github

import (
	"errors"
	"time"
)

var (
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrMissingGitHubToken = errors.New("github token is missing")
)

type ErrorKind string

const (
	ErrorKindRateLimited                ErrorKind = "github_rate_limited"
	ErrorKindTokenInvalid               ErrorKind = "github_token_invalid"
	ErrorKindPermissionDenied           ErrorKind = "github_permission_denied"
	ErrorKindPermissionDeniedOrNotFound ErrorKind = "github_permission_denied_or_not_found"
	ErrorKindUnavailable                ErrorKind = "github_unavailable"
)

type APIError struct {
	Kind       ErrorKind
	Message    string
	StatusCode int
	RetryAfter time.Duration
	ResetAt    *time.Time
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return string(e.Kind)
}
