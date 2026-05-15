// Package guild owns guild concepts used by game use cases.
package guild

import (
	"errors"
	"strings"
	"time"
)

type ID string

type Guild struct {
	ID          ID
	Slug        string
	Name        string
	Description string
	Icon        string
	Color       string
	SortOrder   int
	MemberCount int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewGuild(guild Guild) (Guild, error) {
	if guild.ID == "" {
		return Guild{}, errors.New("guild id is required")
	}
	if strings.TrimSpace(guild.Slug) == "" {
		return Guild{}, errors.New("guild slug is required")
	}
	if strings.TrimSpace(guild.Name) == "" {
		return Guild{}, errors.New("guild name is required")
	}
	if strings.TrimSpace(guild.Description) == "" {
		return Guild{}, errors.New("guild description is required")
	}
	if strings.TrimSpace(guild.Icon) == "" {
		return Guild{}, errors.New("guild icon is required")
	}
	if !isHexColor(guild.Color) {
		return Guild{}, errors.New("guild color must be a 6-digit hex color")
	}
	if guild.SortOrder < 0 {
		return Guild{}, errors.New("guild sort order cannot be negative")
	}
	if guild.MemberCount < 0 {
		return Guild{}, errors.New("guild member count cannot be negative")
	}
	if guild.CreatedAt.IsZero() {
		return Guild{}, errors.New("guild created at is required")
	}
	if guild.UpdatedAt.IsZero() {
		return Guild{}, errors.New("guild updated at is required")
	}

	return guild, nil
}

func isHexColor(value string) bool {
	if len(value) != 7 || value[0] != '#' {
		return false
	}

	for _, char := range value[1:] {
		if !isHexDigit(char) {
			return false
		}
	}

	return true
}

func isHexDigit(char rune) bool {
	return (char >= '0' && char <= '9') ||
		(char >= 'a' && char <= 'f') ||
		(char >= 'A' && char <= 'F')
}
