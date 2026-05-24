INSERT INTO "guild_town_inventories" ("guild_id", "building_type", "quantity", "created_at", "updated_at")
SELECT "id", 'tent', 2, NOW(), NOW()
FROM "guilds"
ON CONFLICT ("guild_id", "building_type")
DO UPDATE SET quantity = EXCLUDED.quantity, updated_at = EXCLUDED.updated_at
WHERE "guild_town_inventories"."quantity" < EXCLUDED.quantity;

INSERT INTO "guild_town_inventories" ("guild_id", "building_type", "quantity", "created_at", "updated_at")
SELECT "id", 'bonfire', 3, NOW(), NOW()
FROM "guilds"
ON CONFLICT ("guild_id", "building_type")
DO UPDATE SET quantity = EXCLUDED.quantity, updated_at = EXCLUDED.updated_at
WHERE "guild_town_inventories"."quantity" < EXCLUDED.quantity;
