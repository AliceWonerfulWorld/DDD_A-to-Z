package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	guildapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guild"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
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

func (s *GuildStore) FindGuildByID(ctx context.Context, guildID guilddomain.ID) (guilddomain.Guild, bool, error) {
	var record guildRecord
	result := s.db.WithContext(ctx).Raw(`
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
		WHERE g.id = ?
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
	`, guildID).Scan(&record)
	if result.Error != nil {
		return guilddomain.Guild{}, false, result.Error
	}
	if result.RowsAffected == 0 {
		return guilddomain.Guild{}, false, nil
	}

	foundGuild, err := record.toDomain()
	if err != nil {
		return guilddomain.Guild{}, false, err
	}

	return foundGuild, true, nil
}

func (s *GuildStore) FindActiveMembershipByUserID(ctx context.Context, userID user.ID) (guilddomain.MembershipWithGuild, bool, error) {
	var record guildMembershipWithGuildRecord
	result := s.db.WithContext(ctx).Raw(`
		SELECT
			gm.id AS membership_id,
			gm.user_id,
			gm.guild_id,
			gm.joined_at,
			gm.left_at,
			gm.created_at AS membership_created_at,
			gm.updated_at AS membership_updated_at,
			g.id,
			g.slug,
			g.name,
			g.description,
			g.icon,
			g.color,
			g.sort_order,
			g.created_at,
			g.updated_at,
			COUNT(active_gm.id) AS member_count
		FROM guild_memberships gm
		JOIN guilds g ON g.id = gm.guild_id
		LEFT JOIN guild_memberships active_gm
			ON active_gm.guild_id = g.id
			AND active_gm.left_at IS NULL
		WHERE gm.user_id = ?
			AND gm.left_at IS NULL
		GROUP BY
			gm.id,
			gm.user_id,
			gm.guild_id,
			gm.joined_at,
			gm.left_at,
			gm.created_at,
			gm.updated_at,
			g.id,
			g.slug,
			g.name,
			g.description,
			g.icon,
			g.color,
			g.sort_order,
			g.created_at,
			g.updated_at
	`, userID).Scan(&record)
	if result.Error != nil {
		return guilddomain.MembershipWithGuild{}, false, result.Error
	}
	if result.RowsAffected == 0 {
		return guilddomain.MembershipWithGuild{}, false, nil
	}

	membership, err := record.toDomain()
	if err != nil {
		return guilddomain.MembershipWithGuild{}, false, err
	}

	return membership, true, nil
}

func (s *GuildStore) CreateMembership(ctx context.Context, membership guilddomain.Membership) error {
	err := s.db.WithContext(ctx).Exec(`
		INSERT INTO guild_memberships (id, user_id, guild_id, joined_at, left_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, membership.ID, membership.UserID, membership.GuildID, membership.JoinedAt, membership.LeftAt, membership.CreatedAt, membership.UpdatedAt).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "guild_memberships_active_user_id_idx" {
			return guildapp.ErrAlreadyJoined
		}
		return err
	}

	return nil
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

type guildMembershipWithGuildRecord struct {
	MembershipID        guilddomain.MembershipID `gorm:"column:membership_id"`
	UserID              user.ID                  `gorm:"column:user_id"`
	GuildID             guilddomain.ID           `gorm:"column:guild_id"`
	JoinedAt            time.Time                `gorm:"column:joined_at"`
	LeftAt              *time.Time               `gorm:"column:left_at"`
	MembershipCreatedAt time.Time                `gorm:"column:membership_created_at"`
	MembershipUpdatedAt time.Time                `gorm:"column:membership_updated_at"`
	ID                  guilddomain.ID           `gorm:"column:id"`
	Slug                string                   `gorm:"column:slug"`
	Name                string                   `gorm:"column:name"`
	Description         string                   `gorm:"column:description"`
	Icon                string                   `gorm:"column:icon"`
	Color               string                   `gorm:"column:color"`
	SortOrder           int                      `gorm:"column:sort_order"`
	MemberCount         int64                    `gorm:"column:member_count"`
	CreatedAt           time.Time                `gorm:"column:created_at"`
	UpdatedAt           time.Time                `gorm:"column:updated_at"`
}

func (r guildMembershipWithGuildRecord) toDomain() (guilddomain.MembershipWithGuild, error) {
	membership, err := guilddomain.NewMembership(guilddomain.Membership{
		ID:        r.MembershipID,
		UserID:    r.UserID,
		GuildID:   r.GuildID,
		JoinedAt:  r.JoinedAt,
		LeftAt:    r.LeftAt,
		CreatedAt: r.MembershipCreatedAt,
		UpdatedAt: r.MembershipUpdatedAt,
	})
	if err != nil {
		return guilddomain.MembershipWithGuild{}, err
	}

	foundGuild, err := guilddomain.NewGuild(guilddomain.Guild{
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
	if err != nil {
		return guilddomain.MembershipWithGuild{}, err
	}

	return guilddomain.MembershipWithGuild{
		Membership: membership,
		Guild:      foundGuild,
	}, nil
}
