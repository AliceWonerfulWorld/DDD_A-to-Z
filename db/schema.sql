CREATE TABLE users (
  id TEXT PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE github_accounts (
  github_id BIGINT PRIMARY KEY,
  user_id TEXT NOT NULL UNIQUE REFERENCES users(id),
  username TEXT NOT NULL,
  avatar_url TEXT NOT NULL,
  access_token_ciphertext TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE sessions (
  token_hash TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_user_id_idx ON sessions(user_id);
CREATE INDEX sessions_expires_at_idx ON sessions(expires_at);

CREATE TABLE point_types (
  code     TEXT NOT NULL,
  language TEXT NOT NULL DEFAULT '',
  label    TEXT NOT NULL,
  PRIMARY KEY (code, language)
);

INSERT INTO point_types (code, language, label) VALUES
  ('CP', '',            'Contribution Point'),
  ('SP', 'Go',          'Go Skill Point'),
  ('SP', 'Python',      'Python Skill Point'),
  ('SP', 'JavaScript',  'JavaScript Skill Point'),
  ('SP', 'TypeScript',  'TypeScript Skill Point'),
  ('SP', 'Rust',        'Rust Skill Point'),
  ('SP', 'Java',        'Java Skill Point'),
  ('SP', 'C',           'C Skill Point'),
  ('SP', 'C++',         'C++ Skill Point'),
  ('SP', 'C#',          'C# Skill Point'),
  ('SP', 'Ruby',        'Ruby Skill Point'),
  ('SP', 'PHP',         'PHP Skill Point'),
  ('SP', 'Swift',       'Swift Skill Point'),
  ('SP', 'Kotlin',      'Kotlin Skill Point'),
  ('SP', 'Scala',       'Scala Skill Point'),
  ('SP', 'Haskell',     'Haskell Skill Point'),
  ('SP', 'Elixir',      'Elixir Skill Point'),
  ('SP', 'Erlang',      'Erlang Skill Point'),
  ('SP', 'Clojure',     'Clojure Skill Point'),
  ('SP', 'Dart',        'Dart Skill Point'),
  ('SP', 'Lua',         'Lua Skill Point'),
  ('SP', 'Shell',       'Shell Skill Point'),
  ('SP', 'PowerShell',  'PowerShell Skill Point'),
  ('SP', 'R',           'R Skill Point'),
  ('SP', 'Julia',       'Julia Skill Point'),
  ('SP', 'Nim',         'Nim Skill Point'),
  ('SP', 'Zig',         'Zig Skill Point'),
  ('SP', 'OCaml',       'OCaml Skill Point'),
  ('SP', 'F#',          'F# Skill Point'),
  ('SP', 'Groovy',      'Groovy Skill Point'),
  ('SP', 'Perl',        'Perl Skill Point'),
  ('SP', 'MATLAB',      'MATLAB Skill Point'),
  ('SP', 'Objective-C', 'Objective-C Skill Point'),
  ('SP', 'Crystal',     'Crystal Skill Point'),
  ('SP', 'Elm',         'Elm Skill Point'),
  ('SP', 'D',           'D Skill Point'),
  ('SP', 'Haxe',        'Haxe Skill Point'),
  ('SP', 'Mojo',        'Mojo Skill Point'),
  ('SP', 'GDScript',    'GDScript Skill Point'),
  ('SP', 'Cuda',        'Cuda Skill Point'),
  ('SP', 'PLpgSQL',     'PLpgSQL Skill Point'),
  ('SP', 'CSS',         'CSS Skill Point'),
  ('SP', 'Nix',         'Nix Skill Point'),
  ('SP', 'HCL',         'HCL Skill Point');

CREATE TABLE point_accounts (
  user_id          TEXT NOT NULL REFERENCES users(id),
  point_type_code  TEXT NOT NULL,
  language         TEXT NOT NULL DEFAULT '',
  balance          BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
  last_analyzed_at TIMESTAMPTZ,
  created_at       TIMESTAMPTZ NOT NULL,
  updated_at       TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (user_id, point_type_code, language),
  FOREIGN KEY (point_type_code, language) REFERENCES point_types(code, language)
);

CREATE TABLE point_ledger (
  id              TEXT PRIMARY KEY,
  user_id         TEXT NOT NULL REFERENCES users(id),
  point_type_code TEXT NOT NULL,
  language        TEXT NOT NULL DEFAULT '',
  amount          BIGINT NOT NULL CHECK (amount <> 0),
  type            TEXT NOT NULL CHECK (type IN ('earn', 'spend', 'adjust')),
  reason          TEXT NOT NULL CHECK (length(reason) > 0),
  source_type     TEXT NOT NULL CHECK (length(source_type) > 0),
  source_id       TEXT NOT NULL CHECK (length(source_id) > 0),
  balance_after   BIGINT NOT NULL CHECK (balance_after >= 0),
  created_at      TIMESTAMPTZ NOT NULL,
  FOREIGN KEY (point_type_code, language) REFERENCES point_types(code, language),
  CHECK (
    (type = 'earn' AND amount > 0)
    OR (type = 'spend' AND amount < 0)
    OR type = 'adjust'
  )
);

CREATE INDEX point_ledger_user_id_created_at_idx ON point_ledger(user_id, point_type_code, language, created_at DESC);
CREATE INDEX point_ledger_source_idx ON point_ledger(source_type, source_id);

CREATE FUNCTION reject_nonzero_point_account_insert()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  IF NEW.balance <> 0 THEN
    RAISE EXCEPTION 'point account must start with zero balance'
      USING ERRCODE = '23514';
  END IF;

  RETURN NEW;
END;
$$;

CREATE TRIGGER point_accounts_reject_nonzero_insert
BEFORE INSERT ON point_accounts
FOR EACH ROW
EXECUTE FUNCTION reject_nonzero_point_account_insert();

CREATE FUNCTION reject_direct_point_account_balance_update()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  IF OLD.balance IS DISTINCT FROM NEW.balance
    AND pg_trigger_depth() < 2 THEN
    RAISE EXCEPTION 'point account balance can only be updated from point_ledger'
      USING ERRCODE = '23514';
  END IF;

  RETURN NEW;
END;
$$;

CREATE TRIGGER point_accounts_reject_direct_balance_update
BEFORE UPDATE OF balance ON point_accounts
FOR EACH ROW
EXECUTE FUNCTION reject_direct_point_account_balance_update();

CREATE FUNCTION apply_point_ledger_entry()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
  next_balance BIGINT;
BEGIN
  SELECT balance + NEW.amount
  INTO next_balance
  FROM point_accounts
  WHERE user_id = NEW.user_id AND point_type_code = NEW.point_type_code AND language = NEW.language
  FOR UPDATE;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'point account not found for user_id % point_type_code % language %', NEW.user_id, NEW.point_type_code, NEW.language
      USING ERRCODE = '23503';
  END IF;

  IF next_balance < 0 THEN
    RAISE EXCEPTION 'point balance cannot be negative for user_id % point_type_code % language %', NEW.user_id, NEW.point_type_code, NEW.language
      USING ERRCODE = '23514';
  END IF;

  UPDATE point_accounts
  SET balance    = next_balance,
      updated_at = NEW.created_at
  WHERE user_id = NEW.user_id AND point_type_code = NEW.point_type_code AND language = NEW.language;

  NEW.balance_after = next_balance;
  RETURN NEW;
END;
$$;

CREATE TRIGGER point_ledger_apply_before_insert
BEFORE INSERT ON point_ledger
FOR EACH ROW
EXECUTE FUNCTION apply_point_ledger_entry();

CREATE FUNCTION reject_point_ledger_mutation()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  RAISE EXCEPTION 'point_ledger is append-only'
    USING ERRCODE = '23514';
END;
$$;

CREATE TRIGGER point_ledger_reject_update
BEFORE UPDATE ON point_ledger
FOR EACH ROW
EXECUTE FUNCTION reject_point_ledger_mutation();

CREATE TRIGGER point_ledger_reject_delete
BEFORE DELETE ON point_ledger
FOR EACH ROW
EXECUTE FUNCTION reject_point_ledger_mutation();

CREATE TABLE guilds (
  id TEXT PRIMARY KEY,
  slug TEXT NOT NULL UNIQUE CHECK (length(slug) > 0),
  name TEXT NOT NULL CHECK (length(name) > 0),
  description TEXT NOT NULL CHECK (length(description) > 0),
  icon TEXT NOT NULL CHECK (length(icon) > 0),
  color TEXT NOT NULL CHECK (color ~ '^#[0-9A-Fa-f]{6}$'),
  sort_order INTEGER NOT NULL CHECK (sort_order >= 0),
  current_exp BIGINT NOT NULL DEFAULT 0 CHECK (current_exp >= 0),
  guild_level INTEGER NOT NULL DEFAULT 1 CHECK (guild_level BETWEEN 1 AND 5),
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

CREATE TABLE guild_cp_contributions (
  id TEXT PRIMARY KEY,
  guild_id TEXT NOT NULL REFERENCES guilds(id),
  user_id TEXT NOT NULL REFERENCES users(id),
  point_ledger_id TEXT NOT NULL UNIQUE REFERENCES point_ledger(id),
  amount BIGINT NOT NULL CHECK (amount > 0),
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX guild_cp_contributions_guild_id_created_at_idx ON guild_cp_contributions(guild_id, created_at DESC);
CREATE INDEX guild_cp_contributions_user_id_created_at_idx ON guild_cp_contributions(user_id, created_at DESC);

CREATE TABLE github_repositories (
  github_id BIGINT NOT NULL,
  user_id TEXT NOT NULL REFERENCES users(id),
  owner TEXT NOT NULL CHECK (length(owner) > 0),
  name TEXT NOT NULL CHECK (length(name) > 0),
  full_name TEXT NOT NULL CHECK (length(full_name) > 0),
  private BOOLEAN NOT NULL,
  fork BOOLEAN NOT NULL,
  archived BOOLEAN NOT NULL,
  default_branch TEXT NOT NULL CHECK (length(default_branch) > 0),
  language TEXT NOT NULL DEFAULT '',
  html_url TEXT NOT NULL CHECK (length(html_url) > 0),
  pushed_at TIMESTAMPTZ,
  github_updated_at TIMESTAMPTZ NOT NULL,
  synced_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (github_id, user_id)
);

CREATE INDEX github_repositories_user_id_pushed_at_idx ON github_repositories(user_id, pushed_at DESC);
CREATE INDEX github_repositories_user_id_full_name_idx ON github_repositories(user_id, full_name);

CREATE TABLE repository_analysis_contributions (
  user_id TEXT NOT NULL REFERENCES users(id),
  repository_full_name TEXT NOT NULL CHECK (length(repository_full_name) > 0),
  contribution_type TEXT NOT NULL CHECK (contribution_type IN ('commit', 'pull_request')),
  external_id TEXT NOT NULL CHECK (length(external_id) > 0),
  message TEXT NOT NULL CHECK (length(message) > 0),
  language TEXT NOT NULL DEFAULT '',
  cp BIGINT NOT NULL CHECK (cp > 0),
  occurred_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (user_id, contribution_type, repository_full_name, external_id)
);

CREATE INDEX repository_analysis_contributions_occurred_at_idx ON repository_analysis_contributions(occurred_at DESC);
CREATE INDEX repository_analysis_contributions_user_id_occurred_at_idx ON repository_analysis_contributions(user_id, occurred_at DESC);

CREATE TABLE badges (
  slug TEXT PRIMARY KEY CHECK (length(slug) > 0),
  name TEXT NOT NULL CHECK (length(name) > 0),
  description TEXT NOT NULL CHECK (length(description) > 0),
  icon TEXT NOT NULL CHECK (length(icon) > 0),
  condition_type TEXT NOT NULL CHECK (length(condition_type) > 0),
  threshold BIGINT NOT NULL CHECK (threshold > 0),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

INSERT INTO badges (slug, name, description, icon, condition_type, threshold, created_at, updated_at) VALUES
  ('cp-bronze',     'CP Bronze',     'Earned 50 Contribution Points',              '🥉', 'cp_earned', 50,  NOW(), NOW()),
  ('cp-silver',     'CP Silver',     'Earned 100 Contribution Points',             '🥈', 'cp_earned', 100, NOW(), NOW()),
  ('cp-gold',       'CP Gold',        'Earned 300 Contribution Points',             '🥇', 'cp_earned', 300, NOW(), NOW()),
  ('cp-platinum',   'CP Platinum',    'Earned 500 Contribution Points',             '💎', 'cp_earned', 500, NOW(), NOW());

CREATE TABLE user_profiles (
  user_id TEXT PRIMARY KEY REFERENCES users(id),
  display_name TEXT NOT NULL CHECK (length(display_name) > 0 AND length(display_name) <= 50),
  selected_badge_slug TEXT REFERENCES badges(slug),
  avatar_url TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE guild_town_inventories (
  guild_id TEXT NOT NULL REFERENCES guilds(id),
  building_type TEXT NOT NULL,
  quantity INTEGER NOT NULL CHECK (quantity >= 0),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (guild_id, building_type)
);

CREATE TABLE guild_town_placements (
  id TEXT PRIMARY KEY,
  guild_id TEXT NOT NULL REFERENCES guilds(id),
  building_type TEXT NOT NULL,
  level INTEGER NOT NULL DEFAULT 1 CHECK (level BETWEEN 1 AND 5),
  x DOUBLE PRECISION NOT NULL CHECK (x >= 0),
  y DOUBLE PRECISION NOT NULL CHECK (y >= 0),
  width DOUBLE PRECISION NOT NULL CHECK (width > 0),
  z_index INTEGER NOT NULL CHECK (z_index >= 0),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX guild_town_placements_guild_id_z_index_idx ON guild_town_placements(guild_id, z_index);

CREATE TABLE chat_tokens (
  token_hash TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id),
  guild_id   TEXT NOT NULL REFERENCES guilds(id),
  expires_at TIMESTAMPTZ NOT NULL,
  used_at    TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX chat_tokens_expires_at_idx ON chat_tokens(expires_at);

CREATE TABLE guild_chat_messages (
  id         TEXT PRIMARY KEY,
  guild_id   TEXT NOT NULL REFERENCES guilds(id),
  user_id    TEXT NOT NULL REFERENCES users(id),
  body       TEXT NOT NULL CHECK (length(body) > 0 AND length(body) <= 1000),
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX guild_chat_messages_guild_id_created_at_idx
  ON guild_chat_messages(guild_id, created_at DESC, id DESC);

CREATE TABLE admin_logs (
  id          BIGSERIAL PRIMARY KEY,
  action      TEXT NOT NULL,
  target_type TEXT NOT NULL,
  target_id   TEXT NOT NULL,
  payload     JSONB NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL
);

CREATE INDEX admin_logs_created_at_idx ON admin_logs(created_at DESC);

CREATE TABLE player_pets (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  guild_id TEXT NOT NULL REFERENCES guilds(id),
  attribute TEXT NOT NULL CHECK (length(attribute) > 0),
  vitality INTEGER NOT NULL CHECK (vitality > 0),
  strength INTEGER NOT NULL CHECK (strength > 0),
  agility INTEGER NOT NULL CHECK (agility > 0),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  UNIQUE (user_id, guild_id)
);

CREATE INDEX player_pets_user_id_created_at_idx ON player_pets(user_id, created_at DESC);

CREATE TABLE tech_news_cache (
  id SERIAL PRIMARY KEY,
  slug TEXT NOT NULL,
  title TEXT NOT NULL,
  url TEXT NOT NULL,
  source TEXT NOT NULL DEFAULT '',
  summary TEXT NOT NULL DEFAULT '',
  published_at TIMESTAMPTZ NOT NULL,
  fetched_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX tech_news_cache_slug_fetched_at_idx ON tech_news_cache(slug, fetched_at DESC);

CREATE TABLE user_badges (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  badge_slug TEXT NOT NULL REFERENCES badges(slug),
  earned_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  UNIQUE (user_id, badge_slug)
);

CREATE INDEX user_badges_user_id_earned_at_idx ON user_badges(user_id, earned_at DESC);

CREATE TABLE seasons (
  id TEXT PRIMARY KEY,
  number INTEGER NOT NULL UNIQUE,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

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

INSERT INTO seasons (id, number, starts_at, ends_at, created_at, updated_at) VALUES
  ('season_1', 1, '2026-05-24 00:00:00+09', '2026-08-24 00:00:00+09', NOW(), NOW());
