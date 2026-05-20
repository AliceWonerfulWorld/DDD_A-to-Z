package postgres

import (
	"context"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/repositoryanalysis"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"gorm.io/gorm"
)

type RepositoryStore struct {
	db *gorm.DB
}

func NewRepositoryStore(db *gorm.DB) *RepositoryStore {
	return &RepositoryStore{db: db}
}

func (s *RepositoryStore) UpsertRepositories(ctx context.Context, repositories []repositoryanalysis.Repository) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, repository := range repositories {
			if _, err := repositoryanalysis.NewRepository(repository); err != nil {
				return err
			}

			if err := tx.Exec(`
				INSERT INTO github_repositories (
					github_id,
					user_id,
					owner,
					name,
					full_name,
					private,
					fork,
					archived,
					default_branch,
					language,
					html_url,
					pushed_at,
					github_updated_at,
					synced_at
				)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT (github_id, user_id) DO UPDATE
				SET owner = EXCLUDED.owner,
					name = EXCLUDED.name,
					full_name = EXCLUDED.full_name,
					private = EXCLUDED.private,
					fork = EXCLUDED.fork,
					archived = EXCLUDED.archived,
					default_branch = EXCLUDED.default_branch,
					language = EXCLUDED.language,
					html_url = EXCLUDED.html_url,
					pushed_at = EXCLUDED.pushed_at,
					github_updated_at = EXCLUDED.github_updated_at,
					synced_at = EXCLUDED.synced_at
			`,
				repository.GitHubID,
				repository.UserID,
				repository.Owner,
				repository.Name,
				repository.FullName,
				repository.Private,
				repository.Fork,
				repository.Archived,
				repository.DefaultBranch,
				repository.Language,
				repository.HTMLURL,
				repository.PushedAt,
				repository.GitHubUpdatedAt,
				repository.SyncedAt,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *RepositoryStore) ListRepositories(ctx context.Context, userID user.ID) ([]repositoryanalysis.Repository, error) {
	var records []repositoryRecord
	if err := s.db.WithContext(ctx).Raw(`
		SELECT
			github_id,
			user_id,
			owner,
			name,
			full_name,
			private,
			fork,
			archived,
			default_branch,
			language,
			html_url,
			pushed_at,
			github_updated_at,
			synced_at
		FROM github_repositories
		WHERE user_id = ?
		ORDER BY pushed_at DESC NULLS LAST, github_updated_at DESC, full_name ASC
	`, userID).Scan(&records).Error; err != nil {
		return nil, err
	}

	repositories := make([]repositoryanalysis.Repository, 0, len(records))
	for _, record := range records {
		repository, err := record.toDomain()
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, repository)
	}

	return repositories, nil
}

func (s *RepositoryStore) UpsertAnalysisContributions(ctx context.Context, userID user.ID, contributions []repositoryanalysis.Contribution, recordedAt time.Time) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, contribution := range contributions {
			if contribution.ExternalID == "" {
				continue
			}
			if contribution.CP <= 0 {
				continue
			}

			if err := tx.Exec(`
				INSERT INTO repository_analysis_contributions (
					user_id,
					repository_full_name,
					contribution_type,
					external_id,
					message,
					language,
					cp,
					occurred_at,
					created_at,
					updated_at
				)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT (user_id, contribution_type, repository_full_name, external_id) DO UPDATE
				SET message = EXCLUDED.message,
					language = EXCLUDED.language,
					cp = EXCLUDED.cp,
					occurred_at = EXCLUDED.occurred_at,
					updated_at = EXCLUDED.updated_at
			`,
				userID,
				contribution.Repo,
				contribution.Type,
				contribution.ExternalID,
				contribution.Message,
				contribution.Language,
				contribution.CP,
				contribution.Timestamp,
				recordedAt,
				recordedAt,
			).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

type repositoryRecord struct {
	GitHubID        int64      `gorm:"column:github_id"`
	UserID          user.ID    `gorm:"column:user_id"`
	Owner           string     `gorm:"column:owner"`
	Name            string     `gorm:"column:name"`
	FullName        string     `gorm:"column:full_name"`
	Private         bool       `gorm:"column:private"`
	Fork            bool       `gorm:"column:fork"`
	Archived        bool       `gorm:"column:archived"`
	DefaultBranch   string     `gorm:"column:default_branch"`
	Language        string     `gorm:"column:language"`
	HTMLURL         string     `gorm:"column:html_url"`
	PushedAt        *time.Time `gorm:"column:pushed_at"`
	GitHubUpdatedAt time.Time  `gorm:"column:github_updated_at"`
	SyncedAt        time.Time  `gorm:"column:synced_at"`
}

func (r repositoryRecord) toDomain() (repositoryanalysis.Repository, error) {
	return repositoryanalysis.NewRepository(repositoryanalysis.Repository{
		GitHubID:        r.GitHubID,
		UserID:          r.UserID,
		Owner:           r.Owner,
		Name:            r.Name,
		FullName:        r.FullName,
		Private:         r.Private,
		Fork:            r.Fork,
		Archived:        r.Archived,
		DefaultBranch:   r.DefaultBranch,
		Language:        r.Language,
		HTMLURL:         r.HTMLURL,
		PushedAt:        r.PushedAt,
		GitHubUpdatedAt: r.GitHubUpdatedAt,
		SyncedAt:        r.SyncedAt,
	})
}
