package pet

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	contributionpointapp "github.com/jyogi-web/ddd-a-to-z/services/api/internal/application/contributionpoint"
	contributionpointdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/contributionpoint"
	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	petdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/pet"
)

var (
	ErrUnauthenticated     = errors.New("unauthenticated")
	ErrPetNotFound         = errors.New("pet not found")
	ErrInvalidTrainStat    = errors.New("invalid training stat")
	ErrInsufficientCP      = errors.New("insufficient cp")
	ErrTrainingIDMissing   = errors.New("pet training id generator is required")
	ErrTrainingUnavailable = errors.New("pet training dependencies are unavailable")
	ErrBattleUnavailable   = errors.New("pet battle dependencies are unavailable")
	ErrOpponentPetNotFound = errors.New("opponent pet not found")
	ErrInvalidBattleTarget = errors.New("invalid battle target")
)

// UseCase handles pet data retrieval for the authenticated user.
type UseCase struct {
	current      CurrentUserRepository
	cp           CPBalanceReader
	pets         PetReader
	guild        CurrentGuildReader
	trainingPets PetTrainingRepository
	trainingCP   CPSpender
	trainingIDs  IDGenerator
	transaction  TrainingTransactioner
	battlePets   PetBattleReader
	now          func() time.Time
}

// NewUseCase creates a new pet use case.
func NewUseCase(current CurrentUserRepository, cp CPBalanceReader, pets PetReader, guild CurrentGuildReader) *UseCase {
	return NewUseCaseWithTraining(current, cp, pets, guild, nil, nil, nil, nil)
}

func NewUseCaseWithTraining(
	current CurrentUserRepository,
	cp CPBalanceReader,
	pets PetReader,
	guild CurrentGuildReader,
	trainingPets PetTrainingRepository,
	trainingCP CPSpender,
	trainingIDs IDGenerator,
	transaction TrainingTransactioner,
) *UseCase {
	return NewUseCaseWithTrainingAndBattle(current, cp, pets, guild, trainingPets, trainingCP, trainingIDs, transaction, nil)
}

func NewUseCaseWithTrainingAndBattle(
	current CurrentUserRepository,
	cp CPBalanceReader,
	pets PetReader,
	guild CurrentGuildReader,
	trainingPets PetTrainingRepository,
	trainingCP CPSpender,
	trainingIDs IDGenerator,
	transaction TrainingTransactioner,
	battlePets PetBattleReader,
) *UseCase {
	return &UseCase{
		current:      current,
		cp:           cp,
		pets:         pets,
		guild:        guild,
		trainingPets: trainingPets,
		trainingCP:   trainingCP,
		trainingIDs:  trainingIDs,
		transaction:  transaction,
		battlePets:   battlePets,
		now:          time.Now,
	}
}

// GetMyPets returns the current user's pets, current guild pet, and CP balance.
func (u *UseCase) GetMyPets(ctx context.Context, sessionToken string) (MyPetsData, error) {
	if sessionToken == "" {
		return MyPetsData{}, ErrUnauthenticated
	}

	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, u.now())
	if err != nil {
		return MyPetsData{}, err
	}
	if !ok {
		return MyPetsData{}, ErrUnauthenticated
	}

	balance, err := u.cp.GetBalance(ctx, appUser.ID)
	if err != nil {
		return MyPetsData{}, err
	}

	pets, err := u.pets.ListPetsByUser(ctx, appUser.ID)
	if err != nil {
		return MyPetsData{}, err
	}

	var currentGuildID guilddomain.ID
	if u.guild != nil {
		membership, found, err := u.guild.FindActiveMembershipByUserID(ctx, appUser.ID)
		if err != nil {
			return MyPetsData{}, err
		}
		if found {
			currentGuildID = membership.Guild.ID
		}
	}

	summaries := make([]PetSummary, 0, len(pets))
	var currentGuildPet *PetSummary
	for _, petWithGuild := range pets {
		summary := toPetSummary(petWithGuild)
		if currentGuildID != "" && petWithGuild.Pet.GuildID == currentGuildID {
			copy := summary
			currentGuildPet = &copy
		}
		summaries = append(summaries, summary)
	}

	return MyPetsData{
		CPBalance:       balance,
		CurrentGuildPet: currentGuildPet,
		Pets:            summaries,
	}, nil
}

func (u *UseCase) TrainPet(ctx context.Context, command TrainPetCommand) (TrainPetResult, error) {
	if err := ctx.Err(); err != nil {
		return TrainPetResult{}, err
	}
	if strings.TrimSpace(command.SessionToken) == "" {
		return TrainPetResult{}, ErrUnauthenticated
	}
	if strings.TrimSpace(command.PetID) == "" {
		return TrainPetResult{}, ErrPetNotFound
	}
	stat, err := petdomain.ParseTrainingStat(command.Stat)
	if err != nil {
		return TrainPetResult{}, ErrInvalidTrainStat
	}
	cost := trainingCost(stat)

	appUser, ok, err := u.current.FindUserBySessionToken(ctx, command.SessionToken, u.now())
	if err != nil {
		return TrainPetResult{}, err
	}
	if !ok {
		return TrainPetResult{}, ErrUnauthenticated
	}
	if u.trainingIDs == nil {
		return TrainPetResult{}, ErrTrainingIDMissing
	}
	trainingID, err := u.trainingIDs.NewID()
	if err != nil {
		return TrainPetResult{}, err
	}

	var result TrainPetResult
	err = u.withTrainingTransaction(ctx, func(ctx context.Context, pets PetTrainingRepository, cp CPSpender) error {
		petWithGuild, found, err := pets.FindPetByIDForUser(ctx, petdomain.ID(command.PetID), appUser.ID)
		if err != nil {
			return err
		}
		if !found {
			return ErrPetNotFound
		}

		ledger, err := cp.Spend(ctx, contributionpointapp.SpendCommand{
			UserID:     appUser.ID,
			PointType:  contributionpointdomain.PointTypeCP,
			Amount:     cost,
			Reason:     fmt.Sprintf("pet_training_%s", stat),
			SourceType: "pet_training",
			SourceID:   trainingID,
		})
		if err != nil {
			if errors.Is(err, contributionpointapp.ErrInsufficientBalance) {
				return ErrInsufficientCP
			}
			return err
		}

		trainedPet, err := petWithGuild.Pet.Train(stat, ledger.CreatedAt)
		if err != nil {
			return err
		}
		if err := pets.UpdatePet(ctx, trainedPet); err != nil {
			return err
		}
		petWithGuild.Pet = trainedPet
		result = TrainPetResult{
			Pet:       toPetSummary(petWithGuild),
			SpentCP:   cost,
			CPBalance: ledger.BalanceAfter,
		}
		return nil
	})
	if err != nil {
		return TrainPetResult{}, err
	}

	return result, nil
}

func (u *UseCase) ListBattleOpponents(ctx context.Context, sessionToken string) (BattleOpponentsData, error) {
	if err := ctx.Err(); err != nil {
		return BattleOpponentsData{}, err
	}
	if strings.TrimSpace(sessionToken) == "" {
		return BattleOpponentsData{}, ErrUnauthenticated
	}
	if u.battlePets == nil {
		return BattleOpponentsData{}, ErrBattleUnavailable
	}

	appUser, ok, err := u.current.FindUserBySessionToken(ctx, sessionToken, u.now())
	if err != nil {
		return BattleOpponentsData{}, err
	}
	if !ok {
		return BattleOpponentsData{}, ErrUnauthenticated
	}

	opponents, err := u.battlePets.ListOpponentPets(ctx, appUser.ID)
	if err != nil {
		return BattleOpponentsData{}, err
	}

	summaries := make([]OpponentSummary, 0, len(opponents))
	for _, opponent := range opponents {
		summaries = append(summaries, toOpponentSummary(opponent))
	}
	return BattleOpponentsData{Opponents: summaries}, nil
}

func (u *UseCase) BattlePet(ctx context.Context, command BattlePetCommand) (BattleResult, error) {
	if err := ctx.Err(); err != nil {
		return BattleResult{}, err
	}
	if strings.TrimSpace(command.SessionToken) == "" {
		return BattleResult{}, ErrUnauthenticated
	}
	if strings.TrimSpace(command.PetID) == "" {
		return BattleResult{}, ErrPetNotFound
	}
	if strings.TrimSpace(command.OpponentPetID) == "" {
		return BattleResult{}, ErrOpponentPetNotFound
	}
	if command.PetID == command.OpponentPetID {
		return BattleResult{}, ErrInvalidBattleTarget
	}
	if u.battlePets == nil {
		return BattleResult{}, ErrBattleUnavailable
	}

	appUser, ok, err := u.current.FindUserBySessionToken(ctx, command.SessionToken, u.now())
	if err != nil {
		return BattleResult{}, err
	}
	if !ok {
		return BattleResult{}, ErrUnauthenticated
	}

	attackerPet, found, err := u.battlePets.FindPetByIDForUser(ctx, petdomain.ID(command.PetID), appUser.ID)
	if err != nil {
		return BattleResult{}, err
	}
	if !found {
		return BattleResult{}, ErrPetNotFound
	}

	defenderPet, found, err := u.battlePets.FindOpponentPetByID(ctx, petdomain.ID(command.OpponentPetID), appUser.ID)
	if err != nil {
		return BattleResult{}, err
	}
	if !found {
		return BattleResult{}, ErrOpponentPetNotFound
	}

	attacker := toPetSummary(attackerPet)
	defender := toPetSummary(defenderPet)
	battle, err := petdomain.Battle(toBattlePet(attacker), toBattlePet(defender))
	if err != nil {
		return BattleResult{}, err
	}

	return toBattleResult(battle, attacker, defender), nil
}

func (u *UseCase) withTrainingTransaction(
	ctx context.Context,
	run func(ctx context.Context, pets PetTrainingRepository, cp CPSpender) error,
) error {
	if u.transaction != nil {
		return u.transaction.WithinPetTraining(ctx, run)
	}
	if u.trainingPets == nil || u.trainingCP == nil {
		return ErrTrainingUnavailable
	}
	return run(ctx, u.trainingPets, u.trainingCP)
}

func trainingCost(stat petdomain.TrainingStat) int64 {
	if stat == petdomain.TrainingStatHP {
		return 20
	}
	return 10
}

func toPetSummary(petWithGuild PetWithGuild) PetSummary {
	foundPet := petWithGuild.Pet
	return PetSummary{
		ID:          string(foundPet.ID),
		OwnerUserID: string(foundPet.UserID),
		GuildID:     string(foundPet.GuildID),
		GuildName:   petWithGuild.Guild.Name,
		Name:        petName(foundPet.Attribute),
		Species:     petSpecies(foundPet.Attribute),
		Attribute:   petAttributeLabel(foundPet.Attribute),
		Level:       1,
		Exp:         0,
		MaxHP:       foundPet.Stats.Vitality*5 + 5,
		Power:       foundPet.Stats.Strength - 1,
		Guard:       foundPet.Stats.Vitality - 1,
		Speed:       foundPet.Stats.Agility - 1,
		AcquiredAt:  foundPet.CreatedAt,
	}
}

func toOpponentSummary(petWithGuild PetWithGuild) OpponentSummary {
	summary := toPetSummary(petWithGuild)
	return OpponentSummary{
		ID:        summary.ID,
		GuildID:   summary.GuildID,
		GuildName: summary.GuildName,
		Name:      summary.Name,
		Species:   summary.Species,
		Attribute: summary.Attribute,
		Level:     summary.Level,
		MaxHP:     summary.MaxHP,
		Power:     summary.Power,
		Guard:     summary.Guard,
		Speed:     summary.Speed,
	}
}

func toBattlePet(summary PetSummary) petdomain.BattlePet {
	return petdomain.BattlePet{
		ID:    petdomain.ID(summary.ID),
		MaxHP: summary.MaxHP,
		Power: summary.Power,
		Guard: summary.Guard,
		Speed: summary.Speed,
	}
}

func toBattleResult(result petdomain.BattleResult, attacker PetSummary, defender PetSummary) BattleResult {
	turns := make([]BattleTurn, 0, len(result.Turns))
	for _, turn := range result.Turns {
		actorName := battlePetName(turn.ActorPetID, attacker, defender)
		targetName := battlePetName(turn.TargetPetID, attacker, defender)
		turns = append(turns, BattleTurn{
			Turn:              turn.Turn,
			ActorPetID:        string(turn.ActorPetID),
			TargetPetID:       string(turn.TargetPetID),
			Damage:            turn.Damage,
			TargetRemainingHP: turn.TargetRemainingHP,
			Message:           fmt.Sprintf("%s attacks %s for %d damage.", actorName, targetName, turn.Damage),
		})
	}

	return BattleResult{
		Result:      battleResultLabel(result.Outcome),
		WinnerPetID: string(result.WinnerPetID),
		Turns:       turns,
		Attacker: BattlePetStatus{
			PetID:       attacker.ID,
			Name:        attacker.Name,
			RemainingHP: result.AttackerRemainingHP,
		},
		Defender: BattlePetStatus{
			PetID:       defender.ID,
			Name:        defender.Name,
			RemainingHP: result.DefenderRemainingHP,
		},
	}
}

func battleResultLabel(outcome petdomain.BattleOutcome) string {
	switch outcome {
	case petdomain.BattleOutcomeAttackerWin:
		return "win"
	case petdomain.BattleOutcomeDefenderWin:
		return "loss"
	default:
		return "draw"
	}
}

func battlePetName(petID petdomain.ID, attacker PetSummary, defender PetSummary) string {
	if petID == petdomain.ID(attacker.ID) {
		return attacker.Name
	}
	if petID == petdomain.ID(defender.ID) {
		return defender.Name
	}
	return string(petID)
}

func petName(attribute petdomain.Attribute) string {
	switch attribute {
	case petdomain.AttributeRust:
		return "Ferris"
	case petdomain.AttributePython:
		return "Py"
	case petdomain.AttributeGo:
		return "Gopher"
	case petdomain.AttributeJava:
		return "Duke"
	case petdomain.AttributeTypeScript:
		return "Scriptie"
	case petdomain.AttributeHaskell:
		return "Lambda"
	case petdomain.AttributeZig:
		return "Ziggy"
	default:
		return string(attribute)
	}
}

func petSpecies(attribute petdomain.Attribute) string {
	switch attribute {
	case petdomain.AttributeRust:
		return "crab"
	case petdomain.AttributePython:
		return "python"
	case petdomain.AttributeGo:
		return "gopher"
	case petdomain.AttributeJava:
		return "duke"
	case petdomain.AttributeTypeScript:
		return "typescript"
	case petdomain.AttributeHaskell:
		return "lambda"
	case petdomain.AttributeZig:
		return "zig"
	default:
		return string(attribute)
	}
}

func petAttributeLabel(attribute petdomain.Attribute) string {
	switch attribute {
	case petdomain.AttributeRust:
		return "Rust"
	case petdomain.AttributePython:
		return "Python"
	case petdomain.AttributeGo:
		return "Go"
	case petdomain.AttributeJava:
		return "Java"
	case petdomain.AttributeTypeScript:
		return "TypeScript"
	case petdomain.AttributeHaskell:
		return "Haskell"
	case petdomain.AttributeZig:
		return "Zig"
	default:
		return string(attribute)
	}
}
