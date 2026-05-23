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