package postgres

import (
	"context"
	"errors"

	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	petapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/pet"
	"gorm.io/gorm"
)

type PetTrainingTransactioner struct {
	db          *gorm.DB
	cpLedgerIDs contributionpointapp.IDGenerator
}

func NewPetTrainingTransactioner(db *gorm.DB, cpLedgerIDs contributionpointapp.IDGenerator) *PetTrainingTransactioner {
	return &PetTrainingTransactioner{
		db:          db,
		cpLedgerIDs: cpLedgerIDs,
	}
}

func (t *PetTrainingTransactioner) WithinPetTraining(
	ctx context.Context,
	run func(ctx context.Context, pets petapp.PetTrainingRepository, cp petapp.CPSpender) error,
) error {
	if t.db == nil {
		return errors.New("db is nil")
	}
	if t.cpLedgerIDs == nil {
		return errors.New("contribution point ledger id generator is required")
	}

	return t.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		petStore := NewPetStore(tx)
		cpStore := NewContributionPointStore(tx)
		cpUseCase := contributionpointapp.NewUseCase(cpStore, t.cpLedgerIDs)

		return run(ctx, petStore, cpUseCase)
	})
}
