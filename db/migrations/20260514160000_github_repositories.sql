-- Modify "github_accounts" table
ALTER TABLE "github_accounts" ADD COLUMN "access_token_ciphertext" text NOT NULL DEFAULT '';
-- Create "github_repositories" table
CREATE TABLE "github_repositories" ("github_id" bigint NOT NULL, "user_id" text NOT NULL, "owner" text NOT NULL, "name" text NOT NULL, "full_name" text NOT NULL, "private" boolean NOT NULL, "fork" boolean NOT NULL, "archived" boolean NOT NULL, "default_branch" text NOT NULL, "language" text NOT NULL DEFAULT '', "html_url" text NOT NULL, "pushed_at" timestamptz NULL, "github_updated_at" timestamptz NOT NULL, "synced_at" timestamptz NOT NULL, PRIMARY KEY ("github_id", "user_id"), CONSTRAINT "github_repositories_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "github_repositories_default_branch_check" CHECK (length(default_branch) > 0), CONSTRAINT "github_repositories_full_name_check" CHECK (length(full_name) > 0), CONSTRAINT "github_repositories_html_url_check" CHECK (length(html_url) > 0), CONSTRAINT "github_repositories_name_check" CHECK (length(name) > 0), CONSTRAINT "github_repositories_owner_check" CHECK (length(owner) > 0));
-- Create index "github_repositories_user_id_full_name_idx" to table: "github_repositories"
CREATE INDEX "github_repositories_user_id_full_name_idx" ON "github_repositories" ("user_id", "full_name");
-- Create index "github_repositories_user_id_pushed_at_idx" to table: "github_repositories"
CREATE INDEX "github_repositories_user_id_pushed_at_idx" ON "github_repositories" ("user_id", "pushed_at" DESC);
