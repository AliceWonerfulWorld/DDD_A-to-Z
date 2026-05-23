package pet

import "errors"

const MaxBattleTurns = 10

type BattleOutcome string

const (
	BattleOutcomeAttackerWin BattleOutcome = "attacker_win"
	BattleOutcomeDefenderWin BattleOutcome = "defender_win"
	BattleOutcomeDraw        BattleOutcome = "draw"
)

type BattlePet struct {
	ID    ID
	MaxHP int
	Power int
	Guard int
	Speed int
}

type BattleTurn struct {
	Turn              int
	ActorPetID        ID
	TargetPetID       ID
	Damage            int
	TargetRemainingHP int
}

type BattleResult struct {
	Outcome             BattleOutcome
	WinnerPetID         ID
	Turns               []BattleTurn
	AttackerRemainingHP int
	DefenderRemainingHP int
}

func Battle(attacker BattlePet, defender BattlePet) (BattleResult, error) {
	if err := validateBattlePet(attacker); err != nil {
		return BattleResult{}, err
	}
	if err := validateBattlePet(defender); err != nil {
		return BattleResult{}, err
	}
	if attacker.ID == defender.ID {
		return BattleResult{}, errors.New("battle pets must be different")
	}

	attackerHP := attacker.MaxHP
	defenderHP := defender.MaxHP
	first, second := attacker, defender
	if defender.Speed > attacker.Speed {
		first, second = defender, attacker
	}

	turns := make([]BattleTurn, 0, MaxBattleTurns)
	for turn := 1; turn <= MaxBattleTurns; turn++ {
		actor := first
		target := second
		if turn%2 == 0 {
			actor = second
			target = first
		}

		damage := Damage(actor, target)
		if target.ID == attacker.ID {
			attackerHP = max(0, attackerHP-damage)
			turns = append(turns, BattleTurn{
				Turn:              turn,
				ActorPetID:        actor.ID,
				TargetPetID:       target.ID,
				Damage:            damage,
				TargetRemainingHP: attackerHP,
			})
			if attackerHP == 0 {
				return BattleResult{
					Outcome:             BattleOutcomeDefenderWin,
					WinnerPetID:         defender.ID,
					Turns:               turns,
					AttackerRemainingHP: attackerHP,
					DefenderRemainingHP: defenderHP,
				}, nil
			}
			continue
		}

		defenderHP = max(0, defenderHP-damage)
		turns = append(turns, BattleTurn{
			Turn:              turn,
			ActorPetID:        actor.ID,
			TargetPetID:       target.ID,
			Damage:            damage,
			TargetRemainingHP: defenderHP,
		})
		if defenderHP == 0 {
			return BattleResult{
				Outcome:             BattleOutcomeAttackerWin,
				WinnerPetID:         attacker.ID,
				Turns:               turns,
				AttackerRemainingHP: attackerHP,
				DefenderRemainingHP: defenderHP,
			}, nil
		}
	}

	outcome := BattleOutcomeDraw
	var winner ID
	if attackerHP > defenderHP {
		outcome = BattleOutcomeAttackerWin
		winner = attacker.ID
	} else if defenderHP > attackerHP {
		outcome = BattleOutcomeDefenderWin
		winner = defender.ID
	}

	return BattleResult{
		Outcome:             outcome,
		WinnerPetID:         winner,
		Turns:               turns,
		AttackerRemainingHP: attackerHP,
		DefenderRemainingHP: defenderHP,
	}, nil
}

func Damage(attacker BattlePet, defender BattlePet) int {
	damage := attacker.Power - defender.Guard/2
	return max(1, damage)
}

func validateBattlePet(pet BattlePet) error {
	if pet.ID == "" {
		return errors.New("battle pet id is required")
	}
	if pet.MaxHP <= 0 {
		return errors.New("battle pet max hp must be positive")
	}
	if pet.Power <= 0 {
		return errors.New("battle pet power must be positive")
	}
	if pet.Guard < 0 {
		return errors.New("battle pet guard must not be negative")
	}
	if pet.Speed < 0 {
		return errors.New("battle pet speed must not be negative")
	}
	return nil
}
