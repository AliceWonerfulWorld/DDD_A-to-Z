package guild

const (
	MaxGuildLevel             = 5
	BuyBuildingExperience     = 300
	GuildTownUpgradeLevel2Exp = 100
	GuildTownUpgradeLevel3Exp = 150
	GuildTownUpgradeLevel4Exp = 200
	GuildTownUpgradeLevel5Exp = 500
)

var guildLevelExperienceThresholds = []int64{0, 5000, 20000, 60000, 150000}

type LevelProgress struct {
	Level                  int
	Experience             int64
	CurrentLevelExperience int64
	NextLevelExperience    int64
}

func GuildLevelProgressFromExperience(experience int64) LevelProgress {
	if experience < 0 {
		experience = 0
	}

	level := 1
	for index, threshold := range guildLevelExperienceThresholds {
		if experience >= threshold {
			level = index + 1
		}
	}
	if level > MaxGuildLevel {
		level = MaxGuildLevel
	}

	current := guildLevelExperienceThresholds[level-1]
	next := current
	if level < MaxGuildLevel {
		next = guildLevelExperienceThresholds[level]
	}

	return LevelProgress{
		Level:                  level,
		Experience:             experience,
		CurrentLevelExperience: current,
		NextLevelExperience:    next,
	}
}

func CalculateUpgradeExp(nextLevel int) int64 {
	switch nextLevel {
	case 2:
		return GuildTownUpgradeLevel2Exp
	case 3:
		return GuildTownUpgradeLevel3Exp
	case 4:
		return GuildTownUpgradeLevel4Exp
	case 5:
		return GuildTownUpgradeLevel5Exp
	default:
		return 0
	}
}
