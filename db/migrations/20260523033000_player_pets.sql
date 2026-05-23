-- Create "player_pets" table
CREATE TABLE "player_pets" (
  "id" text NOT NULL,
  "user_id" text NOT NULL,
  "guild_id" text NOT NULL,
  "attribute" text NOT NULL,
  "vitality" integer NOT NULL,
  "strength" integer NOT NULL,
  "agility" integer NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "player_pets_user_id_guild_id_key" UNIQUE ("user_id", "guild_id"),
  CONSTRAINT "player_pets_guild_id_fkey" FOREIGN KEY ("guild_id") REFERENCES "guilds" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "player_pets_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "player_pets_attribute_check" CHECK (length(attribute) > 0),
  CONSTRAINT "player_pets_vitality_check" CHECK (vitality > 0),
  CONSTRAINT "player_pets_strength_check" CHECK (strength > 0),
  CONSTRAINT "player_pets_agility_check" CHECK (agility > 0)
);
-- Create index "player_pets_user_id_created_at_idx" to table: "player_pets"
CREATE INDEX "player_pets_user_id_created_at_idx" ON "player_pets" ("user_id", "created_at" DESC);
