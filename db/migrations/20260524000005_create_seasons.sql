CREATE TABLE seasons (
  id TEXT PRIMARY KEY,
  number INTEGER NOT NULL UNIQUE,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ NOT NULL,
  is_current BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX seasons_is_current_idx ON seasons(is_current) WHERE is_current;

CREATE TABLE guild_season_rankings (
  id TEXT PRIMARY KEY,
  season_id TEXT NOT NULL REFERENCES seasons(id),
  guild_id TEXT NOT NULL REFERENCES guilds(id),
  total_cp BIGINT NOT NULL DEFAULT 0,
  rank INTEGER NOT NULL,
  member_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  UNIQUE (season_id, guild_id)
);

CREATE INDEX guild_season_rankings_season_rank_idx ON guild_season_rankings(season_id, rank);

CREATE TABLE guild_season_member_rankings (
  id TEXT PRIMARY KEY,
  season_id TEXT NOT NULL REFERENCES seasons(id),
  guild_id TEXT NOT NULL REFERENCES guilds(id),
  user_id TEXT NOT NULL REFERENCES users(id),
  user_name TEXT NOT NULL,
  contributed_cp BIGINT NOT NULL DEFAULT 0,
  rank INTEGER NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  UNIQUE (season_id, guild_id, user_id)
);

CREATE INDEX guild_season_member_rankings_season_guild_rank_idx ON guild_season_member_rankings(season_id, guild_id, rank);