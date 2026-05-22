-- Persist guild level progress on the guild aggregate.
ALTER TABLE "guilds"
ADD COLUMN "current_exp" bigint NOT NULL DEFAULT 0,
ADD COLUMN "guild_level" integer NOT NULL DEFAULT 1,
ADD CONSTRAINT "guilds_current_exp_check" CHECK ("current_exp" >= 0),
ADD CONSTRAINT "guilds_guild_level_check" CHECK ("guild_level" BETWEEN 1 AND 5);

-- Building upgrades need a persisted placement level.
ALTER TABLE "guild_town_placements"
ADD COLUMN "level" integer NOT NULL DEFAULT 1,
ADD CONSTRAINT "guild_town_placements_level_check" CHECK ("level" BETWEEN 1 AND 5);
