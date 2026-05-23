package pet

import (
	"testing"
	"time"

	guilddomain "github.com/jyogi-web/ddd-a-to-z/services/api/internal/domain/guild"
)

func TestNewPet(t *testing.T) {
	now := time.Date(2026, 5, 23, 3, 30, 0, 0, time.UTC)
	valid := Pet{
		ID:        "pet_1",
		UserID:    "user_1",
		GuildID:   "guild_go",
		Attribute: AttributeGo,
		Stats:     Stats{Vitality: 6, Strength: 7, Agility: 7},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := NewPet(valid); err != nil {
		t.Fatalf("NewPet() がエラーを返しました: %v", err)
	}

	tests := []struct {
		name string
		pet  Pet
	}{
		{name: "id が必須", pet: func() Pet {
			pet := valid
			pet.ID = ""
			return pet
		}()},
		{name: "user id が必須", pet: func() Pet {
			pet := valid
			pet.UserID = ""
			return pet
		}()},
		{name: "guild id が必須", pet: func() Pet {
			pet := valid
			pet.GuildID = ""
			return pet
		}()},
		{name: "attribute が必須", pet: func() Pet {
			pet := valid
			pet.Attribute = " "
			return pet
		}()},
		{name: "vitality は正数", pet: func() Pet {
			pet := valid
			pet.Stats.Vitality = 0
			return pet
		}()},
		{name: "strength は正数", pet: func() Pet {
			pet := valid
			pet.Stats.Strength = 0
			return pet
		}()},
		{name: "agility は正数", pet: func() Pet {
			pet := valid
			pet.Stats.Agility = 0
			return pet
		}()},
		{name: "created at が必須", pet: func() Pet {
			pet := valid
			pet.CreatedAt = time.Time{}
			return pet
		}()},
		{name: "updated at が必須", pet: func() Pet {
			pet := valid
			pet.UpdatedAt = time.Time{}
			return pet
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewPet(tt.pet); err == nil {
				t.Fatal("NewPet() error = nil, 期待値 エラー")
			}
		})
	}
}

func TestInitialProfileForGuild(t *testing.T) {
	tests := []struct {
		guildID       guilddomain.ID
		wantAttribute Attribute
	}{
		{guildID: "guild_rust", wantAttribute: AttributeRust},
		{guildID: "guild_python", wantAttribute: AttributePython},
		{guildID: "guild_go", wantAttribute: AttributeGo},
		{guildID: "guild_java", wantAttribute: AttributeJava},
		{guildID: "guild_typescript", wantAttribute: AttributeTypeScript},
		{guildID: "guild_haskell", wantAttribute: AttributeHaskell},
		{guildID: "guild_zig", wantAttribute: AttributeZig},
	}

	for _, tt := range tests {
		t.Run(string(tt.guildID), func(t *testing.T) {
			profile, ok := InitialProfileForGuild(tt.guildID)
			if !ok {
				t.Fatal("InitialProfileForGuild() ok = false, 期待値 true")
			}
			if profile.Attribute != tt.wantAttribute {
				t.Fatalf("Attribute = %q, 期待値 %q", profile.Attribute, tt.wantAttribute)
			}
			if err := validateStats(profile.Stats); err != nil {
				t.Fatalf("initial stats が不正です: %v", err)
			}
		})
	}
}

func TestInitialProfileForGuildReturnsFalseForUnknownGuild(t *testing.T) {
	if _, ok := InitialProfileForGuild("guild_unknown"); ok {
		t.Fatal("InitialProfileForGuild() ok = true, 期待値 false")
	}
}

func TestNewPetFromGuild(t *testing.T) {
	now := time.Date(2026, 5, 23, 3, 30, 0, 0, time.UTC)
	pet, err := NewPetFromGuild("pet_1", "user_1", "guild_go", now)
	if err != nil {
		t.Fatalf("NewPetFromGuild() がエラーを返しました: %v", err)
	}
	if pet.Attribute != AttributeGo {
		t.Fatalf("Attribute = %q, 期待値 %q", pet.Attribute, AttributeGo)
	}
	if pet.Stats != (Stats{Vitality: 6, Strength: 7, Agility: 7}) {
		t.Fatalf("Stats = %+v, 期待値 Go の初期ステータス", pet.Stats)
	}
	if !pet.CreatedAt.Equal(now) || !pet.UpdatedAt.Equal(now) {
		t.Fatalf("timestamps = %v/%v, 期待値 %v", pet.CreatedAt, pet.UpdatedAt, now)
	}
}

func TestNewPetFromGuildRejectsUnknownGuild(t *testing.T) {
	_, err := NewPetFromGuild("pet_1", "user_1", "guild_unknown", time.Now())
	if err == nil {
		t.Fatal("NewPetFromGuild() error = nil, 期待値 エラー")
	}
}
