-- Allow new guild town building IDs without changing DB constraints.
ALTER TABLE "guild_town_inventories"
DROP CONSTRAINT IF EXISTS "guild_town_inventories_building_type_check";

ALTER TABLE "guild_town_placements"
DROP CONSTRAINT IF EXISTS "guild_town_placements_building_type_check";
