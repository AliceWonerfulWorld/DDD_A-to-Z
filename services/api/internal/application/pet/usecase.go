package pet

import (
	"context"
	"errors"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	petdomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/pet"
)

var ErrUnauthenticated = errors.New("unauthenticated")

// UseCase handles pet data retrieval for the authenticated user.
type UseCase struct {
	current CurrentUserRepository
	cp      CPBalanceReader
	pets    PetReader
	guild   CurrentGuildReader
	now     func() time.Time
}

// NewUseCase creates a new pet use case.
func NewUseCase(current CurrentUserRepository, cp CPBalanceReader, pets PetReader, guild CurrentGuildReader) *UseCase {
	return &UseCase{
		current: current,
		cp:      cp,
		pets:    pets,
		guild:   guild,
		now:     time.Now,
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

func toPetSummary(petWithGuild PetWithGuild) PetSummary {
	foundPet := petWithGuild.Pet
	return PetSummary{
		ID:         string(foundPet.ID),
		GuildID:    string(foundPet.GuildID),
		GuildName:  petWithGuild.Guild.Name,
		Name:       petName(foundPet.Attribute),
		Species:    petSpecies(foundPet.Attribute),
		Attribute:  petAttributeLabel(foundPet.Attribute),
		Level:      1,
		Exp:        0,
		MaxHP:      foundPet.Stats.Vitality*5 + 5,
		Power:      foundPet.Stats.Strength - 1,
		Guard:      foundPet.Stats.Vitality - 1,
		Speed:      foundPet.Stats.Agility - 1,
		AcquiredAt: foundPet.CreatedAt,
	}
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
