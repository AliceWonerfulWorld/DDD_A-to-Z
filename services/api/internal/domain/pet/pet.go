// Package pet owns player pet concepts and guild-derived starting stats.
package pet

import (
	"errors"
	"strings"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
	"github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/user"
)

type ID string

type Attribute string

const (
	AttributeRust       Attribute = "rust"
	AttributePython     Attribute = "python"
	AttributeGo         Attribute = "go"
	AttributeJava       Attribute = "java"
	AttributeTypeScript Attribute = "typescript"
	AttributeHaskell    Attribute = "haskell"
	AttributeZig        Attribute = "zig"
)

type Stats struct {
	Vitality int
	Strength int
	Agility  int
}

type Pet struct {
	ID        ID
	UserID    user.ID
	GuildID   guilddomain.ID
	Attribute Attribute
	Stats     Stats
	CreatedAt time.Time
	UpdatedAt time.Time
}

type InitialProfile struct {
	Attribute Attribute
	Stats     Stats
}

var initialProfilesByGuild = map[guilddomain.ID]InitialProfile{
	"guild_rust":       {Attribute: AttributeRust, Stats: Stats{Vitality: 8, Strength: 8, Agility: 4}},
	"guild_python":     {Attribute: AttributePython, Stats: Stats{Vitality: 7, Strength: 6, Agility: 7}},
	"guild_go":         {Attribute: AttributeGo, Stats: Stats{Vitality: 6, Strength: 7, Agility: 7}},
	"guild_java":       {Attribute: AttributeJava, Stats: Stats{Vitality: 9, Strength: 7, Agility: 4}},
	"guild_typescript": {Attribute: AttributeTypeScript, Stats: Stats{Vitality: 5, Strength: 7, Agility: 8}},
	"guild_haskell":    {Attribute: AttributeHaskell, Stats: Stats{Vitality: 6, Strength: 9, Agility: 5}},
	"guild_zig":        {Attribute: AttributeZig, Stats: Stats{Vitality: 5, Strength: 8, Agility: 7}},
}

func NewPet(pet Pet) (Pet, error) {
	if pet.ID == "" {
		return Pet{}, errors.New("pet id is required")
	}
	if pet.UserID == "" {
		return Pet{}, errors.New("pet user id is required")
	}
	if pet.GuildID == "" {
		return Pet{}, errors.New("pet guild id is required")
	}
	if strings.TrimSpace(string(pet.Attribute)) == "" {
		return Pet{}, errors.New("pet attribute is required")
	}
	if err := validateStats(pet.Stats); err != nil {
		return Pet{}, err
	}
	if pet.CreatedAt.IsZero() {
		return Pet{}, errors.New("pet created at is required")
	}
	if pet.UpdatedAt.IsZero() {
		return Pet{}, errors.New("pet updated at is required")
	}

	return pet, nil
}

func NewPetFromGuild(id ID, userID user.ID, guildID guilddomain.ID, now time.Time) (Pet, error) {
	profile, ok := InitialProfileForGuild(guildID)
	if !ok {
		return Pet{}, errors.New("pet initial profile for guild is not defined")
	}

	return NewPet(Pet{
		ID:        id,
		UserID:    userID,
		GuildID:   guildID,
		Attribute: profile.Attribute,
		Stats:     profile.Stats,
		CreatedAt: now,
		UpdatedAt: now,
	})
}

func InitialProfileForGuild(guildID guilddomain.ID) (InitialProfile, bool) {
	profile, ok := initialProfilesByGuild[guildID]
	return profile, ok
}

func validateStats(stats Stats) error {
	if stats.Vitality <= 0 {
		return errors.New("pet vitality must be positive")
	}
	if stats.Strength <= 0 {
		return errors.New("pet strength must be positive")
	}
	if stats.Agility <= 0 {
		return errors.New("pet agility must be positive")
	}

	return nil
}
