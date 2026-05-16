CREATE TABLE guilds (
  id TEXT PRIMARY KEY,
  slug TEXT NOT NULL UNIQUE CHECK (length(slug) > 0),
  name TEXT NOT NULL CHECK (length(name) > 0),
  description TEXT NOT NULL CHECK (length(description) > 0),
  icon TEXT NOT NULL CHECK (length(icon) > 0),
  color TEXT NOT NULL CHECK (color ~ '^#[0-9A-Fa-f]{6}$'),
  sort_order INTEGER NOT NULL CHECK (sort_order >= 0),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX guilds_sort_order_name_idx ON guilds(sort_order, name);

CREATE TABLE guild_memberships (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  guild_id TEXT NOT NULL REFERENCES guilds(id),
  joined_at TIMESTAMPTZ NOT NULL,
  left_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT guild_memberships_left_at_check CHECK (left_at IS NULL OR left_at >= joined_at)
);

CREATE UNIQUE INDEX guild_memberships_active_user_id_idx ON guild_memberships(user_id) WHERE left_at IS NULL;
CREATE INDEX guild_memberships_guild_id_active_idx ON guild_memberships(guild_id) WHERE left_at IS NULL;
