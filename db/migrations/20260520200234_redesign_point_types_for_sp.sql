-- point_accounts の外部キーと PK を先に落とす
-- atlas:nolint CD101
ALTER TABLE "point_accounts" DROP CONSTRAINT "point_accounts_point_type_fkey", DROP CONSTRAINT "point_accounts_pkey";
-- point_ledger の外部キーを先に落とす
-- atlas:nolint CD101
ALTER TABLE "point_ledger" DROP CONSTRAINT "point_ledger_point_type_fkey";
-- Modify "point_types" table
ALTER TABLE "point_types" DROP CONSTRAINT "point_types_pkey", ADD COLUMN "language" text NOT NULL DEFAULT '', ADD PRIMARY KEY ("code", "language");
-- マスターデータを新構造に移行（既存の Golang_SP, TypeScript_SP を新形式に変換）
UPDATE "point_types" SET "code" = 'SP', "language" = 'Go' WHERE "code" = 'Golang_SP';
UPDATE "point_types" SET "code" = 'SP', "language" = 'TypeScript' WHERE "code" = 'TypeScript_SP';
-- CP行に language = '' を設定（DEFAULT '' で既に '' のはずだが明示的に）
UPDATE "point_types" SET "language" = '' WHERE "code" = 'CP';
-- 新しい言語の SP を追加
INSERT INTO "point_types" ("code", "language", "label") VALUES
  ('SP', 'Python',      'Python Skill Point'),
  ('SP', 'JavaScript',  'JavaScript Skill Point'),
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
  ('SP', 'Cuda',        'Cuda Skill Point');
-- Modify "point_accounts" table
-- atlas:nolint DS103
ALTER TABLE "point_accounts" DROP COLUMN "point_type", ADD COLUMN "point_type_code" text NOT NULL DEFAULT 'CP', ADD COLUMN "language" text NOT NULL DEFAULT '', ADD PRIMARY KEY ("user_id", "point_type_code", "language"), ADD CONSTRAINT "point_accounts_point_type_code_language_fkey" FOREIGN KEY ("point_type_code", "language") REFERENCES "point_types" ("code", "language") ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "point_accounts" ALTER COLUMN "point_type_code" DROP DEFAULT;
-- Drop index "point_ledger_user_id_created_at_idx" from table: "point_ledger"
DROP INDEX "point_ledger_user_id_created_at_idx";
-- Modify "point_ledger" table
-- atlas:nolint DS103
ALTER TABLE "point_ledger" DROP COLUMN "point_type", ADD COLUMN "point_type_code" text NOT NULL DEFAULT 'CP', ADD COLUMN "language" text NOT NULL DEFAULT '', ADD CONSTRAINT "point_ledger_point_type_code_language_fkey" FOREIGN KEY ("point_type_code", "language") REFERENCES "point_types" ("code", "language") ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "point_ledger" ALTER COLUMN "point_type_code" DROP DEFAULT;
-- Create index "point_ledger_user_id_created_at_idx" to table: "point_ledger"
CREATE INDEX "point_ledger_user_id_created_at_idx" ON "point_ledger" ("user_id", "point_type_code", "language", "created_at" DESC);
