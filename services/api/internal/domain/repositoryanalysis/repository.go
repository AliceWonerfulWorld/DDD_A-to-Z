// Package repositoryanalysis owns repository concepts used by analysis use cases.
package repositoryanalysis

import (
	"errors"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type Repository struct {
	GitHubID        int64
	UserID          user.ID
	Owner           string
	Name            string
	FullName        string
	Private         bool
	Fork            bool
	Archived        bool
	DefaultBranch   string
	Language        string
	HTMLURL         string
	PushedAt        *time.Time
	GitHubUpdatedAt time.Time
	SyncedAt        time.Time
}

func NewRepository(repository Repository) (Repository, error) {
	if repository.GitHubID <= 0 {
		return Repository{}, errors.New("github repository id must be a positive integer")
	}
	if repository.UserID == "" {
		return Repository{}, errors.New("user id is required")
	}
	if repository.Owner == "" {
		return Repository{}, errors.New("repository owner is required")
	}
	if repository.Name == "" {
		return Repository{}, errors.New("repository name is required")
	}
	if repository.FullName == "" {
		repository.FullName = repository.Owner + "/" + repository.Name
	}
	if repository.DefaultBranch == "" {
		return Repository{}, errors.New("default branch is required")
	}
	if repository.HTMLURL == "" {
		return Repository{}, errors.New("repository html url is required")
	}
	if repository.GitHubUpdatedAt.IsZero() {
		return Repository{}, errors.New("github updated at is required")
	}
	if repository.SyncedAt.IsZero() {
		return Repository{}, errors.New("synced at is required")
	}

	return repository, nil
}
