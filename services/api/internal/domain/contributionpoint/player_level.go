package contributionpoint

import "math"

// PlayerLevelFromTotalEarned returns the player level derived from lifetime CP.
func PlayerLevelFromTotalEarned(totalEarned int64) int {
	if totalEarned <= 0 {
		return 1
	}

	return int(math.Sqrt(float64(totalEarned)/100)) + 1
}

// TotalEarnedForPlayerLevel returns the lifetime CP required to reach a level.
func TotalEarnedForPlayerLevel(level int) int64 {
	if level <= 1 {
		return 0
	}

	levelOffset := int64(level - 1)
	return levelOffset * levelOffset * 100
}

// NextPlayerLevelProgress returns progress toward the next player level.
func NextPlayerLevelProgress(totalEarned int64) (nextLevel int, nextLevelTotalEarned int64, remaining int64) {
	currentLevel := PlayerLevelFromTotalEarned(totalEarned)
	nextLevel = currentLevel + 1
	nextLevelTotalEarned = TotalEarnedForPlayerLevel(nextLevel)
	remaining = nextLevelTotalEarned - totalEarned
	if remaining < 0 {
		remaining = 0
	}

	return nextLevel, nextLevelTotalEarned, remaining
}
