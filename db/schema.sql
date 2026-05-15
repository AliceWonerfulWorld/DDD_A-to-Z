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

CREATE TABLE contribution_point_accounts (
  user_id TEXT PRIMARY KEY REFERENCES users(id),
  balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE contribution_point_ledger (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  amount BIGINT NOT NULL CHECK (amount <> 0),
  type TEXT NOT NULL CHECK (type IN ('earn', 'spend', 'adjust')),
  reason TEXT NOT NULL CHECK (length(reason) > 0),
  source_type TEXT NOT NULL CHECK (length(source_type) > 0),
  source_id TEXT NOT NULL CHECK (length(source_id) > 0),
  balance_after BIGINT NOT NULL CHECK (balance_after >= 0),
  created_at TIMESTAMPTZ NOT NULL,
  CHECK (
    (type = 'earn' AND amount > 0)
    OR (type = 'spend' AND amount < 0)
    OR type = 'adjust'
  )
);

CREATE INDEX contribution_point_ledger_user_id_created_at_idx ON contribution_point_ledger(user_id, created_at DESC);
CREATE INDEX contribution_point_ledger_source_idx ON contribution_point_ledger(source_type, source_id);

CREATE FUNCTION reject_nonzero_contribution_point_account_insert()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  IF NEW.balance <> 0 THEN
    RAISE EXCEPTION 'contribution point account must start with zero balance'
      USING ERRCODE = '23514';
  END IF;

  RETURN NEW;
END;
$$;

CREATE TRIGGER contribution_point_accounts_reject_nonzero_insert
BEFORE INSERT ON contribution_point_accounts
FOR EACH ROW
EXECUTE FUNCTION reject_nonzero_contribution_point_account_insert();

CREATE FUNCTION reject_direct_contribution_point_account_balance_update()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  IF OLD.balance IS DISTINCT FROM NEW.balance
    AND pg_trigger_depth() < 2 THEN
    RAISE EXCEPTION 'contribution point account balance can only be updated from contribution_point_ledger'
      USING ERRCODE = '23514';
  END IF;

  RETURN NEW;
END;
$$;

CREATE TRIGGER contribution_point_accounts_reject_direct_balance_update
BEFORE UPDATE OF balance ON contribution_point_accounts
FOR EACH ROW
EXECUTE FUNCTION reject_direct_contribution_point_account_balance_update();

CREATE FUNCTION apply_contribution_point_ledger_entry()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
  next_balance BIGINT;
BEGIN
  SELECT balance + NEW.amount
  INTO next_balance
  FROM contribution_point_accounts
  WHERE user_id = NEW.user_id
  FOR UPDATE;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'contribution point account not found for user_id %', NEW.user_id
      USING ERRCODE = '23503';
  END IF;

  IF next_balance < 0 THEN
    RAISE EXCEPTION 'contribution point balance cannot be negative for user_id %', NEW.user_id
      USING ERRCODE = '23514';
  END IF;

  UPDATE contribution_point_accounts
  SET balance = next_balance,
      updated_at = NEW.created_at
  WHERE user_id = NEW.user_id;

  NEW.balance_after = next_balance;
  RETURN NEW;
END;
$$;

CREATE TRIGGER contribution_point_ledger_apply_before_insert
BEFORE INSERT ON contribution_point_ledger
FOR EACH ROW
EXECUTE FUNCTION apply_contribution_point_ledger_entry();

CREATE FUNCTION reject_contribution_point_ledger_mutation()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  RAISE EXCEPTION 'contribution_point_ledger is append-only'
    USING ERRCODE = '23514';
END;
$$;

CREATE TRIGGER contribution_point_ledger_reject_update
BEFORE UPDATE ON contribution_point_ledger
FOR EACH ROW
EXECUTE FUNCTION reject_contribution_point_ledger_mutation();

CREATE TRIGGER contribution_point_ledger_reject_delete
BEFORE DELETE ON contribution_point_ledger
FOR EACH ROW
EXECUTE FUNCTION reject_contribution_point_ledger_mutation();

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

CREATE TABLE user_profiles (
  user_id TEXT PRIMARY KEY REFERENCES users(id),
  display_name TEXT NOT NULL CHECK (length(display_name) > 0 AND length(display_name) <= 50),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

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
