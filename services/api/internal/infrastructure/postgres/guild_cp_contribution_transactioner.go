package postgres

import (
	"context"
	"errors"

	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	guildapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/guild"
	"gorm.io/gorm"
)

type GuildCPContributionTransactioner struct {
	db          *gorm.DB
	cpLedgerIDs contributionpointapp.IDGenerator
}

func NewGuildCPContributionTransactioner(db *gorm.DB, cpLedgerIDs contributionpointapp.IDGenerator) *GuildCPContributionTransactioner {
	return &GuildCPContributionTransactioner{
		db:          db,
		cpLedgerIDs: cpLedgerIDs,
	}
}

func (t *GuildCPContributionTransactioner) WithinCPContribution(
	ctx context.Context,
	run func(ctx context.Context, repository guildapp.Repository, cp guildapp.CPSpender) error,
) error {
	if t.db == nil {
		return errors.New("db is nil")
	}
	if t.cpLedgerIDs == nil {
		return errors.New("contribution point ledger id generator is required")
	}

	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		guildStore, err := NewGuildStore(tx)
		if err != nil {
			return err
		}
		cpStore := NewContributionPointStore(tx)
		cpUseCase := contributionpointapp.NewUseCase(cpStore, t.cpLedgerIDs)

		return run(ctx, guildStore, cpUseCase)
	})
}
