// Package mypage owns the my-page aggregation use case.
package mypage

import (
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

// MyPageData is the aggregated response model for the my-page endpoint.
// It gathers user profile, contribution-point summary, repository summary,
// and guild information into a single structure.
type MyPageData struct {
	User         user.User
	CP           CPSummary
	Repositories RepositorySummary
	GitHubStats  *GitHubStats
	Guild        *GuildInfo
	Profile      *ProfileInfo
}

// ProfileInfo holds the user's custom profile settings.
type ProfileInfo struct {
	DisplayName string
	AvatarURL   string
}

// CPSummary holds the current balance and lifetime earn/spend totals.
type CPSummary struct {
	Balance     int64
	TotalEarned int64
	TotalSpent  int64
}

// RepositorySummary holds a compact overview of the user's synced repositories.
type RepositorySummary struct {
	TotalCount      int
	LanguageSummary map[string]int
	Recent          []RecentRepository
}

// RecentRepository is a trimmed-down view of a repository for the my-page.
type RecentRepository struct {
	GitHubID int64
	FullName string
	Language string
	HTMLURL  string
	PushedAt *string
}

// GuildInfo holds the user's guild membership details.
type GuildInfo struct {
	ID          string
	Name        string
	Slug        string
	Icon        string
	Color       string
	Description string
	MemberCount int64
	Rank        int
	TotalGuilds int
	CP          int64
}

// GitHubStats represents the engineer status fetched from GitHub.
type GitHubStats struct {
	TotalStars          int
	TotalPRs            int
	TotalIssues         int
	ContributedTo       int
	PublicRepos         int
	GitHubCreatedAt     string
	YearlyCommits       int
	YearlyContributions int
}
