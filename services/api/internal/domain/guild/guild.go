// Package guild owns guild concepts used by game use cases.
package guild

import (
	"errors"
	"strings"
	"time"

	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type ID string

type Guild struct {
	ID                          ID
	Slug                        string
	Name                        string
	Description                 string
	Icon                        string
	Color                       string
	SortOrder                   int
	MemberCount                 int64
	TotalContributedCP          int64
	GuildExperience             int64
	GuildLevel                  int
	CurrentGuildLevelExperience int64
	NextGuildLevelExperience    int64
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
}

type MembershipID string

type Membership struct {
	ID        MembershipID
	UserID    user.ID
	GuildID   ID
	JoinedAt  time.Time
	LeftAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MembershipWithGuild struct {
	Membership Membership
	Guild      Guild
}

type MemberContribution struct {
	UserID        user.ID
	Name          string
	TotalEarnedCP int64
	JoinedAt      time.Time
}

type ActivityLog struct {
	ID         string
	UserID     user.ID
	Player     string
	Type       string
	Repo       string
	Message    string
	Language   string
	CP         int64
	OccurredAt time.Time
}

const (
	ActivityTypeCommit      = "commit"
	ActivityTypePullRequest = "pull_request"
)

type CPContributionID string

type CPContribution struct {
	ID            CPContributionID
	GuildID       ID
	UserID        user.ID
	PointLedgerID string
	Amount        int64
	CreatedAt     time.Time
}

func NewGuild(guild Guild) (Guild, error) {
	if guild.ID == "" {
		return Guild{}, errors.New("guild id is required")
	}
	if strings.TrimSpace(guild.Slug) == "" {
		return Guild{}, errors.New("guild slug is required")
	}
	if strings.TrimSpace(guild.Name) == "" {
		return Guild{}, errors.New("guild name is required")
	}
	if strings.TrimSpace(guild.Description) == "" {
		return Guild{}, errors.New("guild description is required")
	}
	if strings.TrimSpace(guild.Icon) == "" {
		return Guild{}, errors.New("guild icon is required")
	}
	if !isHexColor(guild.Color) {
		return Guild{}, errors.New("guild color must be a 6-digit hex color")
	}
	if guild.SortOrder < 0 {
		return Guild{}, errors.New("guild sort order cannot be negative")
	}
	if guild.MemberCount < 0 {
		return Guild{}, errors.New("guild member count cannot be negative")
	}
	if guild.TotalContributedCP < 0 {
		return Guild{}, errors.New("guild total contributed cp cannot be negative")
	}
	if guild.GuildExperience < 0 {
		return Guild{}, errors.New("guild experience cannot be negative")
	}
	if guild.CreatedAt.IsZero() {
		return Guild{}, errors.New("guild created at is required")
	}
	if guild.UpdatedAt.IsZero() {
		return Guild{}, errors.New("guild updated at is required")
	}

	progress := GuildLevelProgressFromExperience(guild.GuildExperience)
	guild.GuildLevel = progress.Level
	guild.CurrentGuildLevelExperience = progress.CurrentLevelExperience
	guild.NextGuildLevelExperience = progress.NextLevelExperience

	return guild, nil
}

func (g Guild) AddExperience(amount int64, now time.Time) (Guild, error) {
	if amount < 0 {
		return Guild{}, errors.New("guild experience amount cannot be negative")
	}
	if now.IsZero() {
		return Guild{}, errors.New("guild experience updated at is required")
	}

	g.GuildExperience += amount
	g.UpdatedAt = now

	return NewGuild(g)
}

func NewCPContribution(contribution CPContribution) (CPContribution, error) {
	if contribution.ID == "" {
		return CPContribution{}, errors.New("guild cp contribution id is required")
	}
	if contribution.GuildID == "" {
		return CPContribution{}, errors.New("guild cp contribution guild id is required")
	}
	if contribution.UserID == "" {
		return CPContribution{}, errors.New("guild cp contribution user id is required")
	}
	if strings.TrimSpace(contribution.PointLedgerID) == "" {
		return CPContribution{}, errors.New("guild cp contribution point ledger id is required")
	}
	if contribution.Amount <= 0 {
		return CPContribution{}, errors.New("guild cp contribution amount must be positive")
	}
	if contribution.CreatedAt.IsZero() {
		return CPContribution{}, errors.New("guild cp contribution created at is required")
	}

	return contribution, nil
}

func NewMemberContribution(member MemberContribution) (MemberContribution, error) {
	if member.UserID == "" {
		return MemberContribution{}, errors.New("guild member user id is required")
	}
	if strings.TrimSpace(member.Name) == "" {
		return MemberContribution{}, errors.New("guild member name is required")
	}
	if member.TotalEarnedCP < 0 {
		return MemberContribution{}, errors.New("guild member total earned cp cannot be negative")
	}
	if member.JoinedAt.IsZero() {
		return MemberContribution{}, errors.New("guild member joined at is required")
	}

	return member, nil
}

func NewActivityLog(log ActivityLog) (ActivityLog, error) {
	if strings.TrimSpace(log.ID) == "" {
		return ActivityLog{}, errors.New("guild activity log id is required")
	}
	if log.UserID == "" {
		return ActivityLog{}, errors.New("guild activity log user id is required")
	}
	if strings.TrimSpace(log.Player) == "" {
		return ActivityLog{}, errors.New("guild activity log player is required")
	}
	logType := strings.TrimSpace(log.Type)
	if logType != ActivityTypeCommit && logType != ActivityTypePullRequest {
		return ActivityLog{}, errors.New("invalid guild activity log type: must be 'commit' or 'pull_request'")
	}
	if strings.TrimSpace(log.Repo) == "" {
		return ActivityLog{}, errors.New("guild activity log repo is required")
	}
	if strings.TrimSpace(log.Message) == "" {
		return ActivityLog{}, errors.New("guild activity log message is required")
	}
	if log.CP <= 0 {
		return ActivityLog{}, errors.New("guild activity log cp must be positive")
	}
	if log.OccurredAt.IsZero() {
		return ActivityLog{}, errors.New("guild activity log occurred at is required")
	}

	return log, nil
}

func NewMembership(membership Membership) (Membership, error) {
	if membership.ID == "" {
		return Membership{}, errors.New("guild membership id is required")
	}
	if membership.UserID == "" {
		return Membership{}, errors.New("guild membership user id is required")
	}
	if membership.GuildID == "" {
		return Membership{}, errors.New("guild membership guild id is required")
	}
	if membership.JoinedAt.IsZero() {
		return Membership{}, errors.New("guild membership joined at is required")
	}
	if membership.LeftAt != nil && membership.LeftAt.Before(membership.JoinedAt) {
		return Membership{}, errors.New("guild membership left at cannot be before joined at")
	}
	if membership.CreatedAt.IsZero() {
		return Membership{}, errors.New("guild membership created at is required")
	}
	if membership.UpdatedAt.IsZero() {
		return Membership{}, errors.New("guild membership updated at is required")
	}

	return membership, nil
}

func (m Membership) Active() bool {
	return m.LeftAt == nil
}

func (m Membership) Leave(now time.Time) (Membership, error) {
	if now.IsZero() {
		return Membership{}, errors.New("leave time is required")
	}
	if !m.Active() {
		return Membership{}, errors.New("guild membership is already left")
	}

	leftAt := now
	m.LeftAt = &leftAt
	m.UpdatedAt = now

	return NewMembership(m)
}

func isHexColor(value string) bool {
	if len(value) != 7 || value[0] != '#' {
		return false
	}

	for _, char := range value[1:] {
		if !isHexDigit(char) {
			return false
		}
	}

	return true
}

func isHexDigit(char rune) bool {
	return (char >= '0' && char <= '9') ||
		(char >= 'a' && char <= 'f') ||
		(char >= 'A' && char <= 'F')
}
