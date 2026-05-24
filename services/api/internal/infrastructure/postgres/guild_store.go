package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	guildapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guild"
	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	petdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/pet"
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
			COALESCE(gm.member_count, 0) AS member_count,
			COALESCE(gcc.total_contributed_cp, 0) AS total_contributed_cp,
			g.current_exp AS guild_experience
		FROM guilds g
		LEFT JOIN (
			SELECT guild_id, COUNT(*) AS member_count
			FROM guild_memberships
			WHERE left_at IS NULL
			GROUP BY guild_id
		) gm ON gm.guild_id = g.id
		LEFT JOIN (
			SELECT guild_id, SUM(amount) AS total_contributed_cp
			FROM guild_cp_contributions
			GROUP BY guild_id
		) gcc ON gcc.guild_id = g.id
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
			COALESCE(gm.member_count, 0) AS member_count,
			COALESCE(gcc.total_contributed_cp, 0) AS total_contributed_cp,
			g.current_exp AS guild_experience
		FROM guilds g
		LEFT JOIN (
			SELECT guild_id, COUNT(*) AS member_count
			FROM guild_memberships
			WHERE left_at IS NULL
			GROUP BY guild_id
		) gm ON gm.guild_id = g.id
		LEFT JOIN (
			SELECT guild_id, SUM(amount) AS total_contributed_cp
			FROM guild_cp_contributions
			GROUP BY guild_id
		) gcc ON gcc.guild_id = g.id
		WHERE g.id = ?
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
			COALESCE(active_gm.member_count, 0) AS member_count,
			COALESCE(gcc.total_contributed_cp, 0) AS total_contributed_cp,
			g.current_exp AS guild_experience
		FROM guild_memberships gm
		JOIN guilds g ON g.id = gm.guild_id
		LEFT JOIN (
			SELECT guild_id, COUNT(*) AS member_count
			FROM guild_memberships
			WHERE left_at IS NULL
			GROUP BY guild_id
		) active_gm ON active_gm.guild_id = g.id
		LEFT JOIN (
			SELECT guild_id, SUM(amount) AS total_contributed_cp
			FROM guild_cp_contributions
			GROUP BY guild_id
		) gcc ON gcc.guild_id = g.id
		WHERE gm.user_id = ?
			AND gm.left_at IS NULL
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

func (s *GuildStore) FindMembershipByUserAndGuild(ctx context.Context, userID user.ID, guildID guilddomain.ID) (guilddomain.MembershipWithGuild, bool, error) {
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
			COALESCE(active_gm.member_count, 0) AS member_count,
			COALESCE(gcc.total_contributed_cp, 0) AS total_contributed_cp
		FROM guild_memberships gm
		JOIN guilds g ON g.id = gm.guild_id
		LEFT JOIN (
			SELECT guild_id, COUNT(*) AS member_count
			FROM guild_memberships
			WHERE left_at IS NULL
			GROUP BY guild_id
		) active_gm ON active_gm.guild_id = g.id
		LEFT JOIN (
			SELECT guild_id, SUM(amount) AS total_contributed_cp
			FROM guild_cp_contributions
			GROUP BY guild_id
		) gcc ON gcc.guild_id = g.id
		WHERE gm.user_id = ?
			AND gm.guild_id = ?
			AND gm.left_at IS NULL
	`, userID, guildID).Scan(&record)
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

func (s *GuildStore) ListActiveMembersByGuild(ctx context.Context, guildID guilddomain.ID) ([]guilddomain.MemberContribution, error) {
	var records []guildMemberContributionRecord
	if err := s.db.WithContext(ctx).Raw(`
		SELECT
			gm.user_id,
			COALESCE(up.display_name, ga.username) AS name,
			COALESCE(up.avatar_url, ga.avatar_url) AS avatar_url,
			COALESCE(SUM(pl.amount) FILTER (
				WHERE pl.point_type_code = ? AND pl.language = ? AND pl.type = ?
			), 0) AS total_earned_cp,
			COALESCE(gcc.total_contributed_cp, 0) AS total_contributed_cp,
			gm.joined_at
		FROM guild_memberships gm
		JOIN github_accounts ga ON ga.user_id = gm.user_id
		LEFT JOIN user_profiles up ON up.user_id = gm.user_id
		LEFT JOIN point_ledger pl ON pl.user_id = gm.user_id
		LEFT JOIN (
			SELECT guild_id, user_id, SUM(amount) AS total_contributed_cp
			FROM guild_cp_contributions
			GROUP BY guild_id, user_id
		) gcc ON gcc.user_id = gm.user_id AND gcc.guild_id = gm.guild_id
		WHERE gm.guild_id = ?
			AND gm.left_at IS NULL
		GROUP BY gm.user_id, up.display_name, ga.username, up.avatar_url, ga.avatar_url, gm.joined_at, gcc.total_contributed_cp
		ORDER BY total_contributed_cp DESC, LOWER(COALESCE(up.display_name, ga.username)) ASC, gm.joined_at ASC
	`, contributionpointdomain.PointTypeCP.Code, contributionpointdomain.PointTypeCP.Language, contributionpointdomain.EntryTypeEarn, guildID).Scan(&records).Error; err != nil {
		return nil, err
	}

	members := make([]guilddomain.MemberContribution, 0, len(records))
	for _, record := range records {
		member, err := record.toDomain()
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

func (s *GuildStore) ListActivityLogsByGuild(ctx context.Context, guildID guilddomain.ID, limit int) ([]guilddomain.ActivityLog, error) {
	if limit <= 0 {
		limit = 20
	}

	var records []guildActivityLogRecord
	if err := s.db.WithContext(ctx).Raw(`
		SELECT
			concat(rac.user_id, ':', rac.contribution_type, ':', rac.repository_full_name, ':', rac.external_id) AS id,
			rac.user_id,
			COALESCE(up.display_name, ga.username) AS player,
			rac.contribution_type,
			rac.repository_full_name,
			rac.message,
			rac.language,
			rac.cp,
			rac.occurred_at
		FROM repository_analysis_contributions rac
		JOIN guild_memberships gm ON gm.user_id = rac.user_id
			AND gm.guild_id = ?
			AND gm.left_at IS NULL
			AND rac.occurred_at >= gm.joined_at
		JOIN github_accounts ga ON ga.user_id = rac.user_id
		LEFT JOIN user_profiles up ON up.user_id = rac.user_id
		ORDER BY rac.occurred_at DESC, rac.created_at DESC, rac.external_id DESC
		LIMIT ?
	`, guildID, limit).Scan(&records).Error; err != nil {
		return nil, err
	}

	logs := make([]guilddomain.ActivityLog, 0, len(records))
	for _, record := range records {
		log, err := record.toDomain()
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
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

func (s *GuildStore) UpdateMembership(ctx context.Context, membership guilddomain.Membership) error {
	result := s.db.WithContext(ctx).Exec(`
		UPDATE guild_memberships
		SET left_at = ?, updated_at = ?
		WHERE id = ?
			AND left_at IS NULL
	`, membership.LeftAt, membership.UpdatedAt, membership.ID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return guildapp.ErrActiveMembershipNotFound
	}

	return nil
}

func (s *GuildStore) FindPetByUserAndGuild(ctx context.Context, userID user.ID, guildID guilddomain.ID) (petdomain.Pet, bool, error) {
	var record playerPetRecord
	result := s.db.WithContext(ctx).Raw(`
		SELECT id, user_id, guild_id, attribute, vitality, strength, agility, created_at, updated_at
		FROM player_pets
		WHERE user_id = ?
			AND guild_id = ?
	`, userID, guildID).Scan(&record)
	if result.Error != nil {
		return petdomain.Pet{}, false, result.Error
	}
	if result.RowsAffected == 0 {
		return petdomain.Pet{}, false, nil
	}

	foundPet, err := record.toDomain()
	if err != nil {
		return petdomain.Pet{}, false, err
	}

	return foundPet, true, nil
}

func (s *GuildStore) CreatePet(ctx context.Context, pet petdomain.Pet) error {
	err := s.db.WithContext(ctx).Exec(`
		INSERT INTO player_pets (id, user_id, guild_id, attribute, vitality, strength, agility, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, pet.ID, pet.UserID, pet.GuildID, pet.Attribute, pet.Stats.Vitality, pet.Stats.Strength, pet.Stats.Agility, pet.CreatedAt, pet.UpdatedAt).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "player_pets_user_id_guild_id_key" {
			return guildapp.ErrPetAlreadyOwned
		}
		return err
	}

	return nil
}

func (s *GuildStore) CreateCPContribution(ctx context.Context, contribution guilddomain.CPContribution) error {
	result := s.db.WithContext(ctx).Exec(`
		INSERT INTO guild_cp_contributions (id, guild_id, user_id, point_ledger_id, amount, created_at)
		SELECT ?, ?, ?, pl.id, ?, ?
		FROM point_ledger pl
		WHERE pl.id = ?
			AND pl.user_id = ?
			AND pl.point_type_code = ?
			AND pl.language = ?
			AND pl.type = ?
			AND pl.amount = ?
			AND pl.source_type = ?
			AND pl.source_id = ?
	`,
		contribution.ID,
		contribution.GuildID,
		contribution.UserID,
		contribution.Amount,
		contribution.CreatedAt,
		contribution.PointLedgerID,
		contribution.UserID,
		contributionpointdomain.PointTypeCP.Code,
		contributionpointdomain.PointTypeCP.Language,
		contributionpointdomain.EntryTypeSpend,
		-contribution.Amount,
		"guild_cp_contribution",
		contribution.ID,
	)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return guildapp.ErrInvalidCPContributionLedger
	}

	return nil
}

func (s *GuildStore) ListCPContributionsByGuild(ctx context.Context, guildID guilddomain.ID, limit int) ([]guilddomain.CPContribution, error) {
	return s.listCPContributions(ctx, "guild_id = ?", guildID, limit)
}

func (s *GuildStore) ListCPContributionsByUser(ctx context.Context, userID user.ID, limit int) ([]guilddomain.CPContribution, error) {
	return s.listCPContributions(ctx, "user_id = ?", userID, limit)
}

func (s *GuildStore) listCPContributions(ctx context.Context, where string, value any, limit int) ([]guilddomain.CPContribution, error) {
	if limit <= 0 {
		limit = 50
	}

	var records []guildCPContributionRecord
	if err := s.db.WithContext(ctx).Raw(`
		SELECT id, guild_id, user_id, point_ledger_id, amount, created_at
		FROM guild_cp_contributions
		WHERE `+where+`
		ORDER BY created_at DESC, id DESC
		LIMIT ?
	`, value, limit).Scan(&records).Error; err != nil {
		return nil, err
	}

	contributions := make([]guilddomain.CPContribution, 0, len(records))
	for _, record := range records {
		contribution, err := record.toDomain()
		if err != nil {
			return nil, err
		}
		contributions = append(contributions, contribution)
	}

	return contributions, nil
}

type guildRecord struct {
	ID                 guilddomain.ID `gorm:"column:id"`
	Slug               string         `gorm:"column:slug"`
	Name               string         `gorm:"column:name"`
	Description        string         `gorm:"column:description"`
	Icon               string         `gorm:"column:icon"`
	Color              string         `gorm:"column:color"`
	SortOrder          int            `gorm:"column:sort_order"`
	MemberCount        int64          `gorm:"column:member_count"`
	TotalContributedCP int64          `gorm:"column:total_contributed_cp"`
	GuildExperience    int64          `gorm:"column:guild_experience"`
	CreatedAt          time.Time      `gorm:"column:created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at"`
}

func (r guildRecord) toDomain() (guilddomain.Guild, error) {
	return guilddomain.NewGuild(guilddomain.Guild{
		ID:                 r.ID,
		Slug:               r.Slug,
		Name:               r.Name,
		Description:        r.Description,
		Icon:               r.Icon,
		Color:              r.Color,
		SortOrder:          r.SortOrder,
		MemberCount:        r.MemberCount,
		TotalContributedCP: r.TotalContributedCP,
		GuildExperience:    r.GuildExperience,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
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
	TotalContributedCP  int64                    `gorm:"column:total_contributed_cp"`
	GuildExperience     int64                    `gorm:"column:guild_experience"`
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
		ID:                 r.ID,
		Slug:               r.Slug,
		Name:               r.Name,
		Description:        r.Description,
		Icon:               r.Icon,
		Color:              r.Color,
		SortOrder:          r.SortOrder,
		MemberCount:        r.MemberCount,
		TotalContributedCP: r.TotalContributedCP,
		GuildExperience:    r.GuildExperience,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
	})
	if err != nil {
		return guilddomain.MembershipWithGuild{}, err
	}

	return guilddomain.MembershipWithGuild{
		Membership: membership,
		Guild:      foundGuild,
	}, nil
}

type guildCPContributionRecord struct {
	ID            guilddomain.CPContributionID `gorm:"column:id"`
	GuildID       guilddomain.ID               `gorm:"column:guild_id"`
	UserID        user.ID                      `gorm:"column:user_id"`
	PointLedgerID string                       `gorm:"column:point_ledger_id"`
	Amount        int64                        `gorm:"column:amount"`
	CreatedAt     time.Time                    `gorm:"column:created_at"`
}

type playerPetRecord struct {
	ID        petdomain.ID        `gorm:"column:id"`
	UserID    user.ID             `gorm:"column:user_id"`
	GuildID   guilddomain.ID      `gorm:"column:guild_id"`
	Attribute petdomain.Attribute `gorm:"column:attribute"`
	Vitality  int                 `gorm:"column:vitality"`
	Strength  int                 `gorm:"column:strength"`
	Agility   int                 `gorm:"column:agility"`
	CreatedAt time.Time           `gorm:"column:created_at"`
	UpdatedAt time.Time           `gorm:"column:updated_at"`
}

type guildMemberContributionRecord struct {
	UserID             user.ID   `gorm:"column:user_id"`
	Name               string    `gorm:"column:name"`
	AvatarURL          string    `gorm:"column:avatar_url"`
	TotalEarnedCP      int64     `gorm:"column:total_earned_cp"`
	TotalContributedCP int64     `gorm:"column:total_contributed_cp"`
	JoinedAt           time.Time `gorm:"column:joined_at"`
}

func (r guildMemberContributionRecord) toDomain() (guilddomain.MemberContribution, error) {
	return guilddomain.NewMemberContribution(guilddomain.MemberContribution{
		UserID:             r.UserID,
		Name:               r.Name,
		AvatarURL:          r.AvatarURL,
		TotalEarnedCP:      r.TotalEarnedCP,
		TotalContributedCP: r.TotalContributedCP,
		JoinedAt:           r.JoinedAt,
	})
}

type guildActivityLogRecord struct {
	ID                 string    `gorm:"column:id"`
	UserID             user.ID   `gorm:"column:user_id"`
	Player             string    `gorm:"column:player"`
	ContributionType   string    `gorm:"column:contribution_type"`
	RepositoryFullName string    `gorm:"column:repository_full_name"`
	Message            string    `gorm:"column:message"`
	Language           string    `gorm:"column:language"`
	CP                 int64     `gorm:"column:cp"`
	OccurredAt         time.Time `gorm:"column:occurred_at"`
}

func (r guildActivityLogRecord) toDomain() (guilddomain.ActivityLog, error) {
	return guilddomain.NewActivityLog(guilddomain.ActivityLog{
		ID:         r.ID,
		UserID:     r.UserID,
		Player:     r.Player,
		Type:       r.ContributionType,
		Repo:       r.RepositoryFullName,
		Message:    r.Message,
		Language:   r.Language,
		CP:         r.CP,
		OccurredAt: r.OccurredAt,
	})
}

func (r guildCPContributionRecord) toDomain() (guilddomain.CPContribution, error) {
	return guilddomain.NewCPContribution(guilddomain.CPContribution{
		ID:            r.ID,
		GuildID:       r.GuildID,
		UserID:        r.UserID,
		PointLedgerID: r.PointLedgerID,
		Amount:        r.Amount,
		CreatedAt:     r.CreatedAt,
	})
}

func (r playerPetRecord) toDomain() (petdomain.Pet, error) {
	return petdomain.NewPet(petdomain.Pet{
		ID:        r.ID,
		UserID:    r.UserID,
		GuildID:   r.GuildID,
		Attribute: r.Attribute,
		Stats: petdomain.Stats{
			Vitality: r.Vitality,
			Strength: r.Strength,
			Agility:  r.Agility,
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	})
}
