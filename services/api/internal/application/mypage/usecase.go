package mypage

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

const defaultRecentLimit = 5

var ErrUnauthenticated = errors.New("unauthenticated")

type UseCase struct {
	current      CurrentUserRepository
	cp           ContributionPointReader
	repositories RepositorySummaryReader
	stats        GitHubStatsReader
	tokens       GitHubTokenRepository
	guild        GuildMembershipReader
	profiles     ProfileReader
	now          func() time.Time
}

func NewUseCase(
	current CurrentUserRepository,
	cp ContributionPointReader,
	repositories RepositorySummaryReader,
	stats GitHubStatsReader,
	tokens GitHubTokenRepository,
	guild GuildMembershipReader,
	profiles ProfileReader,
) *UseCase {
	return &UseCase{
		current:      current,
		cp:           cp,
		repositories: repositories,
		stats:        stats,
		tokens:       tokens,
		guild:        guild,
		profiles:     profiles,
		now:          time.Now,
	}
}

func (u *UseCase) GetMyPage(ctx context.Context, sessionToken string) (MyPageData, error) {
	if sessionToken == "" {
		return MyPageData{}, ErrUnauthenticated
	}

	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, u.now())
	if err != nil {
		return MyPageData{}, err
	}
	if !ok {
		return MyPageData{}, ErrUnauthenticated
	}

	balance, err := u.cp.GetBalance(ctx, appUser.ID)
	if err != nil {
		return MyPageData{}, err
	}
	totalEarned, err := u.cp.GetTotalEarned(ctx, appUser.ID)
	if err != nil {
		return MyPageData{}, err
	}
	totalSpent, err := u.cp.GetTotalSpent(ctx, appUser.ID)
	if err != nil {
		return MyPageData{}, err
	}

	repoSummary, err := u.repositories.GetRepositorySummary(ctx, appUser.ID, defaultRecentLimit)
	if err != nil {
		return MyPageData{}, err
	}

	var ghStats *GitHubStats
	accessToken, ok, err := u.tokens.GitHubAccessToken(ctx, appUser.ID)
	if err != nil {
		return MyPageData{}, err
	}
	if ok {
		stats, statsErr := u.stats.FetchStats(ctx, accessToken, appUser.GitHubAccount.Username)
		if statsErr == nil {
			ghStats = stats
		} else {
			slog.WarnContext(ctx, "failed to fetch github stats", "error", statsErr, "username", appUser.GitHubAccount.Username)
		}
	}

	var guildInfo *GuildInfo
	var totalGuilds int
	if u.guild != nil {
		var guildErr error
		guildInfo, guildErr = u.guild.GetGuildMembership(ctx, appUser.ID)
		if guildErr != nil {
			slog.WarnContext(ctx, "failed to get guild membership", "error", guildErr, "user_id", appUser.ID)
		}
		totalGuilds, guildErr = u.guild.GetTotalGuilds(ctx)
		if guildErr != nil {
			slog.WarnContext(ctx, "failed to get total guilds", "error", guildErr, "user_id", appUser.ID)
		}
	}
	if guildInfo != nil {
		guildInfo.TotalGuilds = totalGuilds
	}

	var profileInfo *ProfileInfo
	if u.profiles != nil {
		p, err := u.profiles.GetProfile(ctx, appUser.ID)
		if err == nil {
			profileInfo = p
		}
	}

	return MyPageData{
		User: appUser,
		CP: CPSummary{
			Balance:     balance,
			TotalEarned: totalEarned,
			TotalSpent:  totalSpent,
		},
		Repositories: repoSummary,
		GitHubStats:  ghStats,
		Guild:        guildInfo,
		Profile:      profileInfo,
	}, nil
}
