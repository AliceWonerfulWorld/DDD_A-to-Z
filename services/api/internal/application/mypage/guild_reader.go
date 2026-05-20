package mypage

import (
	"context"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type GuildRepository interface {
	FindActiveMembershipByUserID(ctx context.Context, userID user.ID) (guilddomain.MembershipWithGuild, bool, error)
	ListGuilds(ctx context.Context) ([]guilddomain.Guild, error)
}

type GuildMembershipAdapter struct {
	repo GuildRepository
}

func NewGuildMembershipReader(repo GuildRepository) GuildMembershipReader {
	return &GuildMembershipAdapter{repo: repo}
}

func (r *GuildMembershipAdapter) GetGuildMembership(ctx context.Context, userID user.ID) (*GuildInfo, error) {
	membership, found, err := r.repo.FindActiveMembershipByUserID(ctx, userID)
	if err != nil || !found {
		return nil, nil
	}

	guild := membership.Guild
	return &GuildInfo{
		ID:          string(guild.ID),
		Name:        guild.Name,
		Slug:        guild.Slug,
		Icon:        guild.Icon,
		Color:       guild.Color,
		Description: guild.Description,
		MemberCount: guild.MemberCount,
		Rank:        0,
		TotalGuilds: 0,
		CP:          0,
	}, nil
}

func (r *GuildMembershipAdapter) GetTotalGuilds(ctx context.Context) (int, error) {
	guilds, err := r.repo.ListGuilds(ctx)
	if err != nil {
		return 0, nil
	}
	return len(guilds), nil
}
