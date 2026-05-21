package guild

const (
	MaxGuildLevel                 = 5
	GuildTownPlacementExperience  = 250
	GuildTownUpgradeExperience    = 500
	GuildExperienceSourceTypeLike = "guild_%"
)

var guildLevelExperienceThresholds = []int64{0, 1000, 3000, 7000, 15000}

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
