package postgres

import (
	"context"
	"errors"
	"time"

	petapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/pet"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	petdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/pet"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"gorm.io/gorm"
)

// PetStore provides read-only queries for player pets.
type PetStore struct {
	db *gorm.DB
}

// NewPetStore creates a new PetStore.
func NewPetStore(db *gorm.DB) *PetStore {
	return &PetStore{db: db}
}

// ListPetsByUser returns all pets owned by a user with their guild display data.
func (s *PetStore) ListPetsByUser(ctx context.Context, userID user.ID) ([]petapp.PetWithGuild, error) {
	var records []playerPetWithGuildRecord
	if err := s.db.WithContext(ctx).Raw(`
		SELECT
			pp.id,
			pp.user_id,
			pp.guild_id,
			pp.attribute,
			pp.vitality,
			pp.strength,
			pp.agility,
			pp.created_at,
			pp.updated_at,
			g.slug,
			g.name,
			g.description,
			g.icon,
			g.color,
			g.sort_order,
			g.created_at AS guild_created_at,
			g.updated_at AS guild_updated_at,
			COALESCE(active_gm.member_count, 0) AS member_count,
			COALESCE(gcc.total_contributed_cp, 0) AS total_contributed_cp,
			g.current_exp AS guild_experience
		FROM player_pets pp
		JOIN guilds g ON g.id = pp.guild_id
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
		WHERE pp.user_id = ?
		ORDER BY pp.created_at DESC, pp.id ASC
	`, userID).Scan(&records).Error; err != nil {
		return nil, err
	}

	pets := make([]petapp.PetWithGuild, 0, len(records))
	for _, record := range records {
		petWithGuild, err := record.toApplicationModel()
		if err != nil {
			return nil, err
		}
		pets = append(pets, petWithGuild)
	}

	return pets, nil
}

func (s *PetStore) FindPetByIDForUser(ctx context.Context, petID petdomain.ID, userID user.ID) (petapp.PetWithGuild, bool, error) {
	var record playerPetWithGuildRecord
	result := s.db.WithContext(ctx).Raw(`
		SELECT
			pp.id,
			pp.user_id,
			pp.guild_id,
			pp.attribute,
			pp.vitality,
			pp.strength,
			pp.agility,
			pp.created_at,
			pp.updated_at,
			g.slug,
			g.name,
			g.description,
			g.icon,
			g.color,
			g.sort_order,
			g.created_at AS guild_created_at,
			g.updated_at AS guild_updated_at,
			COALESCE(active_gm.member_count, 0) AS member_count,
			COALESCE(gcc.total_contributed_cp, 0) AS total_contributed_cp,
			g.current_exp AS guild_experience
		FROM player_pets pp
		JOIN guilds g ON g.id = pp.guild_id
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
		WHERE pp.id = ? AND pp.user_id = ?
		FOR UPDATE OF pp
	`, petID, userID).Scan(&record)
	if result.Error != nil {
		return petapp.PetWithGuild{}, false, result.Error
	}
	if result.RowsAffected == 0 {
		return petapp.PetWithGuild{}, false, nil
	}

	petWithGuild, err := record.toApplicationModel()
	if err != nil {
		return petapp.PetWithGuild{}, false, err
	}
	return petWithGuild, true, nil
}

func (s *PetStore) ListOpponentPets(ctx context.Context, userID user.ID) ([]petapp.PetWithGuild, error) {
	var records []playerPetWithGuildRecord
	if err := s.db.WithContext(ctx).Raw(`
		SELECT
			pp.id,
			pp.user_id,
			pp.guild_id,
			pp.attribute,
			pp.vitality,
			pp.strength,
			pp.agility,
			pp.created_at,
			pp.updated_at,
			g.slug,
			g.name,
			g.description,
			g.icon,
			g.color,
			g.sort_order,
			g.created_at AS guild_created_at,
			g.updated_at AS guild_updated_at,
			COALESCE(active_gm.member_count, 0) AS member_count,
			COALESCE(gcc.total_contributed_cp, 0) AS total_contributed_cp,
			g.current_exp AS guild_experience,
			up.display_name
		FROM player_pets pp
		JOIN guilds g ON g.id = pp.guild_id
		JOIN user_profiles up ON up.user_id = pp.user_id
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
		WHERE pp.user_id <> ?
		ORDER BY pp.created_at DESC, pp.id ASC
		LIMIT 20
	`, userID).Scan(&records).Error; err != nil {
		return nil, err
	}

	pets := make([]petapp.PetWithGuild, 0, len(records))
	for _, record := range records {
		petWithGuild, err := record.toApplicationModel()
		if err != nil {
			return nil, err
		}
		pets = append(pets, petWithGuild)
	}

	return pets, nil
}

func (s *PetStore) FindOpponentPetByID(ctx context.Context, petID petdomain.ID, userID user.ID) (petapp.PetWithGuild, bool, error) {
	var record playerPetWithGuildRecord
	result := s.db.WithContext(ctx).Raw(`
		SELECT
			pp.id,
			pp.user_id,
			pp.guild_id,
			pp.attribute,
			pp.vitality,
			pp.strength,
			pp.agility,
			pp.created_at,
			pp.updated_at,
			g.slug,
			g.name,
			g.description,
			g.icon,
			g.color,
			g.sort_order,
			g.created_at AS guild_created_at,
			g.updated_at AS guild_updated_at,
			COALESCE(active_gm.member_count, 0) AS member_count,
			COALESCE(gcc.total_contributed_cp, 0) AS total_contributed_cp,
			g.current_exp AS guild_experience
		FROM player_pets pp
		JOIN guilds g ON g.id = pp.guild_id
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
		WHERE pp.id = ? AND pp.user_id <> ?
	`, petID, userID).Scan(&record)
	if result.Error != nil {
		return petapp.PetWithGuild{}, false, result.Error
	}
	if result.RowsAffected == 0 {
		return petapp.PetWithGuild{}, false, nil
	}

	petWithGuild, err := record.toApplicationModel()
	if err != nil {
		return petapp.PetWithGuild{}, false, err
	}
	return petWithGuild, true, nil
}

func (s *PetStore) UpdatePet(ctx context.Context, pet petdomain.Pet) error {
	result := s.db.WithContext(ctx).
		Table("player_pets").
		Where("id = ? AND user_id = ?", pet.ID, pet.UserID).
		Updates(map[string]any{
			"vitality":   pet.Stats.Vitality,
			"strength":   pet.Stats.Strength,
			"agility":    pet.Stats.Agility,
			"updated_at": pet.UpdatedAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("pet not found")
	}
	return nil
}

type playerPetWithGuildRecord struct {
	ID                 petdomain.ID        `gorm:"column:id"`
	UserID             user.ID             `gorm:"column:user_id"`
	GuildID            guilddomain.ID      `gorm:"column:guild_id"`
	DisplayName        string              `gorm:"column:display_name"`
	Attribute          petdomain.Attribute `gorm:"column:attribute"`
	Vitality           int                 `gorm:"column:vitality"`
	Strength           int                 `gorm:"column:strength"`
	Agility            int                 `gorm:"column:agility"`
	CreatedAt          time.Time           `gorm:"column:created_at"`
	UpdatedAt          time.Time           `gorm:"column:updated_at"`
	Slug               string              `gorm:"column:slug"`
	Name               string              `gorm:"column:name"`
	Description        string              `gorm:"column:description"`
	Icon               string              `gorm:"column:icon"`
	Color              string              `gorm:"column:color"`
	SortOrder          int                 `gorm:"column:sort_order"`
	MemberCount        int64               `gorm:"column:member_count"`
	TotalContributedCP int64               `gorm:"column:total_contributed_cp"`
	GuildExperience    int64               `gorm:"column:guild_experience"`
	GuildCreatedAt     time.Time           `gorm:"column:guild_created_at"`
	GuildUpdatedAt     time.Time           `gorm:"column:guild_updated_at"`
}

func (r playerPetWithGuildRecord) toApplicationModel() (petapp.PetWithGuild, error) {
	foundPet, err := petdomain.NewPet(petdomain.Pet{
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
	if err != nil {
		return petapp.PetWithGuild{}, err
	}

	foundGuild, err := guilddomain.NewGuild(guilddomain.Guild{
		ID:                 r.GuildID,
		Slug:               r.Slug,
		Name:               r.Name,
		Description:        r.Description,
		Icon:               r.Icon,
		Color:              r.Color,
		SortOrder:          r.SortOrder,
		MemberCount:        r.MemberCount,
		TotalContributedCP: r.TotalContributedCP,
		GuildExperience:    r.GuildExperience,
		CreatedAt:          r.GuildCreatedAt,
		UpdatedAt:          r.GuildUpdatedAt,
	})
	if err != nil {
		return petapp.PetWithGuild{}, err
	}

	return petapp.PetWithGuild{Pet: foundPet, Guild: foundGuild, DisplayName: r.DisplayName}, nil
}
