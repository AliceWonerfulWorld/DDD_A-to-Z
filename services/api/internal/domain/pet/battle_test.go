package pet

import "testing"

func TestBattleSpeedDecidesFirstAttack(t *testing.T) {
	attacker := BattlePet{ID: "pet_attacker", MaxHP: 30, Power: 6, Guard: 4, Speed: 3}
	defender := BattlePet{ID: "pet_defender", MaxHP: 30, Power: 6, Guard: 4, Speed: 7}

	result, err := Battle(attacker, defender)
	if err != nil {
		t.Fatalf("Battle() error = %v", err)
	}
	if len(result.Turns) == 0 {
		t.Fatal("turns length = 0")
	}
	if result.Turns[0].ActorPetID != defender.ID {
		t.Fatalf("first actor = %q, 期待値 %q", result.Turns[0].ActorPetID, defender.ID)
	}
}

func TestBattleSameSpeedAttackerActsFirst(t *testing.T) {
	attacker := BattlePet{ID: "pet_attacker", MaxHP: 30, Power: 6, Guard: 4, Speed: 5}
	defender := BattlePet{ID: "pet_defender", MaxHP: 30, Power: 6, Guard: 4, Speed: 5}

	result, err := Battle(attacker, defender)
	if err != nil {
		t.Fatalf("Battle() error = %v", err)
	}
	if result.Turns[0].ActorPetID != attacker.ID {
		t.Fatalf("first actor = %q, 期待値 %q", result.Turns[0].ActorPetID, attacker.ID)
	}
}

func TestDamageUsesPowerAndHalfGuardWithMinimumOne(t *testing.T) {
	damage := Damage(
		BattlePet{ID: "pet_attacker", MaxHP: 30, Power: 5, Guard: 1, Speed: 1},
		BattlePet{ID: "pet_defender", MaxHP: 30, Power: 1, Guard: 8, Speed: 1},
	)
	if damage != 1 {
		t.Fatalf("damage = %d, 期待値 1", damage)
	}

	damage = Damage(
		BattlePet{ID: "pet_attacker", MaxHP: 30, Power: 9, Guard: 1, Speed: 1},
		BattlePet{ID: "pet_defender", MaxHP: 30, Power: 1, Guard: 5, Speed: 1},
	)
	if damage != 7 {
		t.Fatalf("damage = %d, 期待値 7", damage)
	}
}

func TestBattleStopsWhenPetIsDefeated(t *testing.T) {
	attacker := BattlePet{ID: "pet_attacker", MaxHP: 20, Power: 20, Guard: 1, Speed: 10}
	defender := BattlePet{ID: "pet_defender", MaxHP: 10, Power: 1, Guard: 1, Speed: 1}

	result, err := Battle(attacker, defender)
	if err != nil {
		t.Fatalf("Battle() error = %v", err)
	}
	if result.Outcome != BattleOutcomeAttackerWin {
		t.Fatalf("outcome = %q, 期待値 %q", result.Outcome, BattleOutcomeAttackerWin)
	}
	if result.WinnerPetID != attacker.ID {
		t.Fatalf("winner = %q, 期待値 %q", result.WinnerPetID, attacker.ID)
	}
	if len(result.Turns) != 1 {
		t.Fatalf("turns length = %d, 期待値 1", len(result.Turns))
	}
	if result.DefenderRemainingHP != 0 {
		t.Fatalf("defender HP = %d, 期待値 0", result.DefenderRemainingHP)
	}
}

func TestBattleDecidesByRemainingHPAfterTenTurns(t *testing.T) {
	attacker := BattlePet{ID: "pet_attacker", MaxHP: 50, Power: 2, Guard: 6, Speed: 5}
	defender := BattlePet{ID: "pet_defender", MaxHP: 40, Power: 2, Guard: 6, Speed: 4}

	result, err := Battle(attacker, defender)
	if err != nil {
		t.Fatalf("Battle() error = %v", err)
	}
	if len(result.Turns) != MaxBattleTurns {
		t.Fatalf("turns length = %d, 期待値 %d", len(result.Turns), MaxBattleTurns)
	}
	if result.Outcome != BattleOutcomeAttackerWin {
		t.Fatalf("outcome = %q, 期待値 %q", result.Outcome, BattleOutcomeAttackerWin)
	}
	if result.WinnerPetID != attacker.ID {
		t.Fatalf("winner = %q, 期待値 %q", result.WinnerPetID, attacker.ID)
	}
}

func TestBattleDrawWhenRemainingHPIsSameAfterTenTurns(t *testing.T) {
	attacker := BattlePet{ID: "pet_attacker", MaxHP: 50, Power: 2, Guard: 6, Speed: 5}
	defender := BattlePet{ID: "pet_defender", MaxHP: 50, Power: 2, Guard: 6, Speed: 4}

	result, err := Battle(attacker, defender)
	if err != nil {
		t.Fatalf("Battle() error = %v", err)
	}
	if result.Outcome != BattleOutcomeDraw {
		t.Fatalf("outcome = %q, 期待値 %q", result.Outcome, BattleOutcomeDraw)
	}
	if result.WinnerPetID != "" {
		t.Fatalf("winner = %q, 期待値 empty", result.WinnerPetID)
	}
}
