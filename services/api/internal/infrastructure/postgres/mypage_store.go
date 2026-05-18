package postgres

import (
	"context"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/mypage"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"gorm.io/gorm"
)

// MyPageStore provides read-only queries for the my-page aggregation.
type MyPageStore struct {
	db *gorm.DB
}

// NewMyPageStore creates a new MyPageStore.
func NewMyPageStore(db *gorm.DB) *MyPageStore {
	return &MyPageStore{db: db}
}

// GetTotalEarned returns the sum of all positive (earn) ledger entries.
func (s *MyPageStore) GetTotalEarned(ctx context.Context, userID user.ID) (int64, error) {
	var total *int64
	err := s.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(amount), 0)
		FROM point_ledger
		WHERE user_id = ? AND point_type = 'CP' AND type = 'earn'
	`, userID).Scan(&total).Error
	if err != nil {
		return 0, err
	}
	if total == nil {
		return 0, nil
	}

	return *total, nil
}

// GetTotalSpent returns the absolute sum of all negative (spend) ledger entries.
func (s *MyPageStore) GetTotalSpent(ctx context.Context, userID user.ID) (int64, error) {
	var total *int64
	err := s.db.WithContext(ctx).Raw(`
		SELECT COALESCE(ABS(SUM(amount)), 0)
		FROM point_ledger
		WHERE user_id = ? AND point_type = 'CP' AND type = 'spend'
	`, userID).Scan(&total).Error
	if err != nil {
		return 0, err
	}
	if total == nil {
		return 0, nil
	}

	return *total, nil
}

// GetTodayEarned returns the sum of CP earned today.
func (s *MyPageStore) GetTodayEarned(ctx context.Context, userID user.ID) (int64, error) {
	var total *int64
	err := s.db.WithContext(ctx).Raw(`
		SELECT COALESCE(SUM(amount), 0)
		FROM point_ledger
		WHERE user_id = ? AND point_type = 'CP' AND type = 'earn' AND created_at >= CURRENT_DATE
	`, userID).Scan(&total).Error
	if err != nil {
		return 0, err
	}
	if total == nil {
		return 0, nil
	}

	return *total, nil
}

// GetRepositorySummary returns total repository count, language distribution,
// and the most recently pushed repositories.
func (s *MyPageStore) GetRepositorySummary(ctx context.Context, userID user.ID, recentLimit int) (mypage.RepositorySummary, error) {
	// Total count
	var totalCount int64
	if err := s.db.WithContext(ctx).Raw(`
		SELECT COUNT(*) FROM github_repositories WHERE user_id = ?
	`, userID).Scan(&totalCount).Error; err != nil {
		return mypage.RepositorySummary{}, err
	}

	// Language summary: count repositories per language (excluding empty)
	type langCount struct {
		Language string `gorm:"column:language"`
		Count    int    `gorm:"column:count"`
	}
	var langCounts []langCount
	if err := s.db.WithContext(ctx).Raw(`
		SELECT language, COUNT(*) as count
		FROM github_repositories
		WHERE user_id = ? AND language <> ''
		GROUP BY language
		ORDER BY count DESC
	`, userID).Scan(&langCounts).Error; err != nil {
		return mypage.RepositorySummary{}, err
	}
	langSummary := make(map[string]int, len(langCounts))
	for _, lc := range langCounts {
		langSummary[lc.Language] = lc.Count
	}

	// Recent repositories
	type recentRecord struct {
		GitHubID int64      `gorm:"column:github_id"`
		FullName string     `gorm:"column:full_name"`
		Language string     `gorm:"column:language"`
		HTMLURL  string     `gorm:"column:html_url"`
		PushedAt *time.Time `gorm:"column:pushed_at"`
	}
	var recentRecords []recentRecord
	if err := s.db.WithContext(ctx).Raw(`
		SELECT github_id, full_name, language, html_url, pushed_at
		FROM github_repositories
		WHERE user_id = ?
		ORDER BY pushed_at DESC NULLS LAST, full_name ASC
		LIMIT ?
	`, userID, recentLimit).Scan(&recentRecords).Error; err != nil {
		return mypage.RepositorySummary{}, err
	}

	recent := make([]mypage.RecentRepository, 0, len(recentRecords))
	for _, r := range recentRecords {
		repo := mypage.RecentRepository{
			GitHubID: r.GitHubID,
			FullName: r.FullName,
			Language: r.Language,
			HTMLURL:  r.HTMLURL,
		}
		if r.PushedAt != nil {
			t := r.PushedAt.Format(time.RFC3339)
			repo.PushedAt = &t
		}
		recent = append(recent, repo)
	}

	return mypage.RepositorySummary{
		TotalCount:      int(totalCount),
		LanguageSummary: langSummary,
		Recent:          recent,
	}, nil
}
