// Package pet owns the my-pets retrieval use case.
package pet

import "time"

// MyPetsData is the response model for the authenticated user's pets.
type MyPetsData struct {
	CPBalance       int64
	CurrentGuildPet *PetSummary
	Pets            []PetSummary
}

type TrainPetCommand struct {
	SessionToken string
	PetID        string
	Stat         string
}

type TrainPetResult struct {
	Pet       PetSummary
	SpentCP   int64
	CPBalance int64
}

type BattlePetCommand struct {
	SessionToken  string
	PetID         string
	OpponentPetID string
}

type BattleOpponentsData struct {
	Opponents []OpponentSummary
}

type OpponentSummary struct {
	ID        string
	GuildID   string
	GuildName string
	Name      string
	Species   string
	Attribute string
	Level     int
	MaxHP     int
	Power     int
	Guard     int
	Speed     int
}

type BattleResult struct {
	Result      string
	WinnerPetID string
	Turns       []BattleTurn
	Attacker    BattlePetStatus
	Defender    BattlePetStatus
}

type BattleTurn struct {
	Turn              int
	ActorPetID        string
	TargetPetID       string
	Damage            int
	TargetRemainingHP int
	Message           string
}

type BattlePetStatus struct {
	PetID       string
	Name        string
	RemainingHP int
}

// PetSummary is a frontend-oriented view of a player pet.
type PetSummary struct {
	ID          string
	OwnerUserID string
	GuildID     string
	GuildName   string
	Name        string
	Species     string
	Attribute   string
	Level       int
	Exp         int64
	MaxHP       int
	Power       int
	Guard       int
	Speed       int
	AcquiredAt  time.Time
}
