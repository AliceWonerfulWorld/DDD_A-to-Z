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
