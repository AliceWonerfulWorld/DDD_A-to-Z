package postgres

import (
	"context"
	"errors"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	"gorm.io/gorm"
)

type GuildStore struct {
	db *gorm.DB
}

func NewGuildStore(db *gorm.DB) (*GuildStore, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &GuildStore{db: db}, nil
}

func (s *GuildStore) ListGuilds(ctx context.Context) ([]guilddomain.Guild, error) {
	var records []guildRecord
	if err := s.db.WithContext(ctx).Raw(`
		SELECT
			g.id,
			g.slug,
			g.name,
			g.description,
			g.icon,
			g.color,
			g.sort_order,
			g.created_at,
			g.updated_at,
			COUNT(gm.id) AS member_count
		FROM guilds g
		LEFT JOIN guild_memberships gm
			ON gm.guild_id = g.id
			AND gm.left_at IS NULL
		GROUP BY
			g.id,
			g.slug,
			g.name,
			g.description,
			g.icon,
			g.color,
			g.sort_order,
			g.created_at,
			g.updated_at
		ORDER BY g.sort_order ASC, g.name ASC
	`).Scan(&records).Error; err != nil {
		return nil, err
	}

	guilds := make([]guilddomain.Guild, 0, len(records))
	for _, record := range records {
		guild, err := record.toDomain()
		if err != nil {
			return nil, err
		}
		guilds = append(guilds, guild)
	}

	return guilds, nil
}

type guildRecord struct {
	ID          guilddomain.ID `gorm:"column:id"`
	Slug        string         `gorm:"column:slug"`
	Name        string         `gorm:"column:name"`
	Description string         `gorm:"column:description"`
	Icon        string         `gorm:"column:icon"`
	Color       string         `gorm:"column:color"`
	SortOrder   int            `gorm:"column:sort_order"`
	MemberCount int64          `gorm:"column:member_count"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
}

func (r guildRecord) toDomain() (guilddomain.Guild, error) {
	return guilddomain.NewGuild(guilddomain.Guild{
		ID:          r.ID,
		Slug:        r.Slug,
		Name:        r.Name,
		Description: r.Description,
		Icon:        r.Icon,
		Color:       r.Color,
		SortOrder:   r.SortOrder,
		MemberCount: r.MemberCount,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	})
}
