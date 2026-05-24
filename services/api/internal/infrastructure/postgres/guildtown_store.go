package postgres

import (
	"context"
	"time"

	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	guildtownapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guildtown"
	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	guildtowndomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guildtown"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
	"gorm.io/gorm"
)

type GuildTownStore struct {
	db          *gorm.DB
	cpLedgerIDs contributionpointapp.IDGenerator
}

func NewGuildTownStore(db *gorm.DB) *GuildTownStore {
	return &GuildTownStore{db: db}
}

func NewGuildTownStoreWithLedgerIDs(db *gorm.DB, cpLedgerIDs contributionpointapp.IDGenerator) *GuildTownStore {
	return &GuildTownStore{db: db, cpLedgerIDs: cpLedgerIDs}
}

func (s *GuildTownStore) ListInventory(ctx context.Context, guildID guilddomain.ID) ([]guildtowndomain.InventoryItem, error) {
	var records []guildTownInventoryRecord
	if err := s.db.WithContext(ctx).Raw(`
		SELECT guild_id, building_type, quantity, created_at, updated_at
		FROM guild_town_inventories
		WHERE guild_id = ?
		ORDER BY building_type ASC
	`, guildID).Scan(&records).Error; err != nil {
		return nil, err
	}

	items := make([]guildtowndomain.InventoryItem, 0, len(records))
	for _, record := range records {
		item, err := record.toDomain()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (s *GuildTownStore) ListPlacements(ctx context.Context, guildID guilddomain.ID) ([]guildtowndomain.Placement, error) {
	var records []guildTownPlacementRecord
	if err := s.db.WithContext(ctx).Raw(`
		SELECT id, guild_id, building_type, level, x, y, width, z_index, created_at, updated_at
		FROM guild_town_placements
		WHERE guild_id = ?
		ORDER BY z_index ASC, created_at ASC, id ASC
	`, guildID).Scan(&records).Error; err != nil {
		return nil, err
	}

	placements := make([]guildtowndomain.Placement, 0, len(records))
	for _, record := range records {
		placement, err := record.toDomain()
		if err != nil {
			return nil, err
		}
		placements = append(placements, placement)
	}

	return placements, nil
}

func (s *GuildTownStore) FindPlacementByID(ctx context.Context, guildID guilddomain.ID, placementID guildtowndomain.PlacementID) (guildtowndomain.Placement, bool, error) {
	var record guildTownPlacementRecord
	result := s.db.WithContext(ctx).Raw(`
		SELECT id, guild_id, building_type, level, x, y, width, z_index, created_at, updated_at
		FROM guild_town_placements
		WHERE guild_id = ? AND id = ?
	`, guildID, placementID).Scan(&record)
	if result.Error != nil {
		return guildtowndomain.Placement{}, false, result.Error
	}
	if result.RowsAffected == 0 {
		return guildtowndomain.Placement{}, false, nil
	}

	placement, err := record.toDomain()
	if err != nil {
		return guildtowndomain.Placement{}, false, err
	}

	return placement, true, nil
}

func (s *GuildTownStore) BuyBuilding(ctx context.Context, userID user.ID, guildID guilddomain.ID, building guildtowndomain.BuildingMaster, exp int64, now time.Time) (guilddomain.Guild, error) {
	var updatedGuild guilddomain.Guild
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockGuildForUpdate(ctx, tx, guildID); err != nil {
			return err
		}
		if err := s.spendPurchaseCost(ctx, tx, userID, guildID, building); err != nil {
			return err
		}

		if err := tx.Exec(`
			INSERT INTO guild_town_inventories (guild_id, building_type, quantity, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT (guild_id, building_type)
			DO UPDATE SET quantity = guild_town_inventories.quantity + 1, updated_at = EXCLUDED.updated_at
		`, guildID, building.Type, guildtowndomain.DefaultInventoryQuantity(building.Type)+1, now, now).Error; err != nil {
			return err
		}

		guild, err := addGuildExperience(ctx, tx, guildID, exp, now)
		if err != nil {
			return err
		}
		updatedGuild = guild
		return nil
	})
	if err != nil {
		return guilddomain.Guild{}, err
	}

	return updatedGuild, nil
}

func (s *GuildTownStore) spendPurchaseCost(ctx context.Context, tx *gorm.DB, userID user.ID, guildID guilddomain.ID, building guildtowndomain.BuildingMaster) error {
	if building.PurchaseCP == 0 && building.PurchaseSP == 0 {
		return nil
	}
	if s.cpLedgerIDs == nil {
		return nil
	}

	cp := contributionpointapp.NewUseCase(NewContributionPointStore(tx), s.cpLedgerIDs)
	sourceID := string(guildID) + ":" + string(building.Type)
	if building.PurchaseCP > 0 {
		if _, err := cp.Spend(ctx, contributionpointapp.SpendCommand{
			UserID:     userID,
			PointType:  contributionpointdomain.PointTypeCP,
			Amount:     building.PurchaseCP,
			Reason:     "guild_town_building_purchase",
			SourceType: "guild_town_building",
			SourceID:   sourceID,
		}); err != nil {
			return err
		}
	}
	if shouldSpendPurchaseSP(building) {
		if _, err := cp.Spend(ctx, contributionpointapp.SpendCommand{
			UserID:     userID,
			PointType:  contributionpointdomain.SPType(building.TargetSP),
			Amount:     building.PurchaseSP,
			Reason:     "guild_town_building_purchase",
			SourceType: "guild_town_building",
			SourceID:   sourceID,
		}); err != nil {
			return err
		}
	}

	return nil
}

func shouldSpendPurchaseSP(building guildtowndomain.BuildingMaster) bool {
	// Common is a cross-language bucket in the town UI, not a persisted SP account.
	// Until point_types has a canonical SP/Common entry and earning source, purchases only charge CP for Common-targeted buildings.
	return building.PurchaseSP > 0 && building.TargetSP != "" && building.TargetSP != "Common"
}

func (s *GuildTownStore) CreatePlacement(ctx context.Context, guildID guilddomain.ID, placement guildtowndomain.Placement) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockGuildForUpdate(ctx, tx, guildID); err != nil {
			return err
		}

		var quantity int
		if err := tx.Raw(`
			SELECT quantity
			FROM guild_town_inventories
			WHERE guild_id = ? AND building_type = ?
			FOR UPDATE
		`, guildID, placement.BuildingType).Scan(&quantity).Error; err != nil {
			return err
		}
		if defaultQuantity := guildtowndomain.DefaultInventoryQuantity(placement.BuildingType); quantity < defaultQuantity {
			quantity = defaultQuantity
		}

		var placedCount int64
		if err := tx.Raw(`
			SELECT COUNT(*)
			FROM guild_town_placements
			WHERE guild_id = ? AND building_type = ?
		`, guildID, placement.BuildingType).Scan(&placedCount).Error; err != nil {
			return err
		}
		if int(placedCount) >= quantity {
			return guildtownapp.ErrInsufficientInventory
		}

		var nextZIndex int
		if err := tx.Raw(`
			SELECT COALESCE(MAX(z_index) + 1, 0)
			FROM guild_town_placements
			WHERE guild_id = ?
		`, guildID).Scan(&nextZIndex).Error; err != nil {
			return err
		}

		return tx.Exec(`
			INSERT INTO guild_town_placements (
				id, guild_id, building_type, level, x, y, width, z_index, created_at, updated_at
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, placement.ID, guildID, placement.BuildingType, placement.Level, placement.X, placement.Y, placement.Width, nextZIndex, placement.CreatedAt, placement.UpdatedAt).Error
	})
}

func (s *GuildTownStore) ReplacePlacements(ctx context.Context, guildID guilddomain.ID, placements []guildtowndomain.Placement) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockGuildForUpdate(ctx, tx, guildID); err != nil {
			return err
		}

		if err := tx.Exec(`DELETE FROM guild_town_placements WHERE guild_id = ?`, guildID).Error; err != nil {
			return err
		}

		for _, placement := range placements {
			if err := tx.Exec(`
				INSERT INTO guild_town_placements (
					id, guild_id, building_type, level, x, y, width, z_index, created_at, updated_at
				)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, placement.ID, guildID, placement.BuildingType, placement.Level, placement.X, placement.Y, placement.Width, placement.ZIndex, placement.CreatedAt, placement.UpdatedAt).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *GuildTownStore) UpgradePlacement(ctx context.Context, guildID guilddomain.ID, placementID guildtowndomain.PlacementID, nextLevel int, exp int64, now time.Time) (guilddomain.Guild, error) {
	var updatedGuild guilddomain.Guild
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := lockGuildForUpdate(ctx, tx, guildID); err != nil {
			return err
		}

		var currentLevel int
		result := tx.Raw(`
			SELECT level
			FROM guild_town_placements
			WHERE guild_id = ? AND id = ?
			FOR UPDATE
		`, guildID, placementID).Scan(&currentLevel)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return guildtownapp.ErrPlacementNotFound
		}
		if nextLevel != currentLevel+1 || nextLevel > guilddomain.MaxGuildLevel {
			return guildtownapp.ErrInvalidPlacementLevel
		}

		if err := tx.Exec(`
			UPDATE guild_town_placements
			SET level = ?, updated_at = ?
			WHERE guild_id = ? AND id = ?
		`, nextLevel, now, guildID, placementID).Error; err != nil {
			return err
		}

		guild, err := addGuildExperience(ctx, tx, guildID, exp, now)
		if err != nil {
			return err
		}
		updatedGuild = guild
		return nil
	})
	if err != nil {
		return guilddomain.Guild{}, err
	}

	return updatedGuild, nil
}

type guildTownInventoryRecord struct {
	GuildID      guilddomain.ID               `gorm:"column:guild_id"`
	BuildingType guildtowndomain.BuildingType `gorm:"column:building_type"`
	Quantity     int                          `gorm:"column:quantity"`
	CreatedAt    time.Time                    `gorm:"column:created_at"`
	UpdatedAt    time.Time                    `gorm:"column:updated_at"`
}

func (r guildTownInventoryRecord) toDomain() (guildtowndomain.InventoryItem, error) {
	return guildtowndomain.NewInventoryItem(guildtowndomain.InventoryItem{
		GuildID:      r.GuildID,
		BuildingType: r.BuildingType,
		Quantity:     r.Quantity,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	})
}

type guildTownPlacementRecord struct {
	ID           guildtowndomain.PlacementID  `gorm:"column:id"`
	GuildID      guilddomain.ID               `gorm:"column:guild_id"`
	BuildingType guildtowndomain.BuildingType `gorm:"column:building_type"`
	Level        int                          `gorm:"column:level"`
	X            float64                      `gorm:"column:x"`
	Y            float64                      `gorm:"column:y"`
	Width        float64                      `gorm:"column:width"`
	ZIndex       int                          `gorm:"column:z_index"`
	CreatedAt    time.Time                    `gorm:"column:created_at"`
	UpdatedAt    time.Time                    `gorm:"column:updated_at"`
}

func (r guildTownPlacementRecord) toDomain() (guildtowndomain.Placement, error) {
	return guildtowndomain.NewPlacement(guildtowndomain.Placement{
		ID:           r.ID,
		GuildID:      r.GuildID,
		BuildingType: r.BuildingType,
		Level:        r.Level,
		X:            r.X,
		Y:            r.Y,
		Width:        r.Width,
		ZIndex:       r.ZIndex,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	})
}

type guildForUpdateRecord struct {
	Found              bool           `gorm:"column:found"`
	ID                 guilddomain.ID `gorm:"column:id"`
	Slug               string         `gorm:"column:slug"`
	Name               string         `gorm:"column:name"`
	Description        string         `gorm:"column:description"`
	Icon               string         `gorm:"column:icon"`
	Color              string         `gorm:"column:color"`
	SortOrder          int            `gorm:"column:sort_order"`
	MemberCount        int64          `gorm:"column:member_count"`
	TotalContributedCP int64          `gorm:"column:total_contributed_cp"`
	CurrentExp         int64          `gorm:"column:current_exp"`
	CreatedAt          time.Time      `gorm:"column:created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at"`
}

func lockGuildForUpdate(ctx context.Context, tx *gorm.DB, guildID guilddomain.ID) error {
	var found bool
	result := tx.WithContext(ctx).Raw(`
		SELECT TRUE
		FROM guilds
		WHERE id = ?
		FOR UPDATE
	`, guildID).Scan(&found)
	if result.Error != nil {
		return result.Error
	}
	if !found {
		return guildtownapp.ErrGuildNotFound
	}

	return nil
}

func addGuildExperience(ctx context.Context, tx *gorm.DB, guildID guilddomain.ID, exp int64, now time.Time) (guilddomain.Guild, error) {
	var record guildForUpdateRecord
	result := tx.WithContext(ctx).Raw(`
		SELECT
			TRUE AS found,
			g.id,
			g.slug,
			g.name,
			g.description,
			g.icon,
			g.color,
			g.sort_order,
			g.created_at,
			g.updated_at,
			g.current_exp,
			COALESCE(gm.member_count, 0) AS member_count,
			COALESCE(gcc.total_contributed_cp, 0) AS total_contributed_cp
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
		FOR UPDATE OF g
	`, guildID).Scan(&record)
	if result.Error != nil {
		return guilddomain.Guild{}, result.Error
	}
	if !record.Found {
		return guilddomain.Guild{}, guildtownapp.ErrGuildNotFound
	}

	guild, err := guilddomain.NewGuild(guilddomain.Guild{
		ID:                 record.ID,
		Slug:               record.Slug,
		Name:               record.Name,
		Description:        record.Description,
		Icon:               record.Icon,
		Color:              record.Color,
		SortOrder:          record.SortOrder,
		MemberCount:        record.MemberCount,
		TotalContributedCP: record.TotalContributedCP,
		GuildExperience:    record.CurrentExp,
		CreatedAt:          record.CreatedAt,
		UpdatedAt:          record.UpdatedAt,
	})
	if err != nil {
		return guilddomain.Guild{}, err
	}

	updatedGuild, err := guild.AddExperience(exp, now)
	if err != nil {
		return guilddomain.Guild{}, err
	}

	if err := tx.WithContext(ctx).Exec(`
		UPDATE guilds
		SET current_exp = ?, guild_level = ?, updated_at = ?
		WHERE id = ?
	`, updatedGuild.GuildExperience, updatedGuild.GuildLevel, updatedGuild.UpdatedAt, guildID).Error; err != nil {
		return guilddomain.Guild{}, err
	}

	return updatedGuild, nil
}
