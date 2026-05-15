-- Create "guilds" table
CREATE TABLE "guilds" (
  "id" text NOT NULL,
  "slug" text NOT NULL,
  "name" text NOT NULL,
  "description" text NOT NULL,
  "icon" text NOT NULL,
  "color" text NOT NULL,
  "sort_order" integer NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "guilds_slug_key" UNIQUE ("slug"),
  CONSTRAINT "guilds_color_check" CHECK (color ~ '^#[0-9A-Fa-f]{6}$'),
  CONSTRAINT "guilds_description_check" CHECK (length(description) > 0),
  CONSTRAINT "guilds_icon_check" CHECK (length(icon) > 0),
  CONSTRAINT "guilds_name_check" CHECK (length(name) > 0),
  CONSTRAINT "guilds_slug_check" CHECK (length(slug) > 0),
  CONSTRAINT "guilds_sort_order_check" CHECK (sort_order >= 0)
);
-- Create index "guilds_sort_order_name_idx" to table: "guilds"
CREATE INDEX "guilds_sort_order_name_idx" ON "guilds" ("sort_order", "name");
-- Seed MVP guilds
INSERT INTO "guilds" ("id", "slug", "name", "description", "icon", "color", "sort_order", "created_at", "updated_at")
VALUES
  ('guild_rust', 'rust', 'Rust', '安全性と速度を掲げる、堅牢なシステム開発者のギルド。', '🦀', '#ff6b35', 1, '2026-05-15T00:00:00Z', '2026-05-15T00:00:00Z'),
  ('guild_python', 'python', 'Python', '学習しやすさと汎用性で課題を解く、万能型のギルド。', '🐍', '#3776ab', 2, '2026-05-15T00:00:00Z', '2026-05-15T00:00:00Z'),
  ('guild_go', 'go', 'Go', 'シンプルさと並列処理で前に進む、軽快なバックエンドギルド。', '🐹', '#00acd7', 3, '2026-05-15T00:00:00Z', '2026-05-15T00:00:00Z'),
  ('guild_java', 'java', 'Java', '長い歴史と堅牢なエコシステムを誇る、重装騎士のギルド。', '☕', '#f89820', 4, '2026-05-15T00:00:00Z', '2026-05-15T00:00:00Z'),
  ('guild_typescript', 'typescript', 'TypeScript', '型の力でフロントからバックエンドまで支える、適応型のギルド。', '📘', '#3178c6', 5, '2026-05-15T00:00:00Z', '2026-05-15T00:00:00Z'),
  ('guild_haskell', 'haskell', 'Haskell', '理論と型安全性を突き詰める、静かな賢者たちのギルド。', 'λ', '#5d4f85', 6, '2026-05-15T00:00:00Z', '2026-05-15T00:00:00Z'),
  ('guild_zig', 'zig', 'Zig', '低レイヤーをまっすぐ扱う、原始の力を持つギルド。', '⚡', '#f7a41d', 7, '2026-05-15T00:00:00Z', '2026-05-15T00:00:00Z');
-- Create "guild_memberships" table
CREATE TABLE "guild_memberships" (
  "id" text NOT NULL,
  "user_id" text NOT NULL,
  "guild_id" text NOT NULL,
  "joined_at" timestamptz NOT NULL,
  "left_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "guild_memberships_guild_id_fkey" FOREIGN KEY ("guild_id") REFERENCES "guilds" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "guild_memberships_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "guild_memberships_left_at_check" CHECK (left_at IS NULL OR left_at >= joined_at)
);
-- Create index "guild_memberships_active_user_id_idx" to table: "guild_memberships"
CREATE UNIQUE INDEX "guild_memberships_active_user_id_idx" ON "guild_memberships" ("user_id") WHERE left_at IS NULL;
-- Create index "guild_memberships_guild_id_active_idx" to table: "guild_memberships"
CREATE INDEX "guild_memberships_guild_id_active_idx" ON "guild_memberships" ("guild_id") WHERE left_at IS NULL;
