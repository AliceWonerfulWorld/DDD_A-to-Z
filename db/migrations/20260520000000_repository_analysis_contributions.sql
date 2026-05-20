-- Create "repository_analysis_contributions" table
CREATE TABLE "repository_analysis_contributions" (
  "user_id" text NOT NULL,
  "repository_full_name" text NOT NULL,
  "contribution_type" text NOT NULL,
  "external_id" text NOT NULL,
  "message" text NOT NULL,
  "language" text NOT NULL DEFAULT '',
  "cp" bigint NOT NULL,
  "occurred_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("user_id", "contribution_type", "repository_full_name", "external_id"),
  CONSTRAINT "repository_analysis_contributions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "repository_analysis_contributions_contribution_type_check" CHECK (contribution_type IN ('commit', 'pull_request')),
  CONSTRAINT "repository_analysis_contributions_cp_check" CHECK (cp > 0),
  CONSTRAINT "repository_analysis_contributions_external_id_check" CHECK (length(external_id) > 0),
  CONSTRAINT "repository_analysis_contributions_message_check" CHECK (length(message) > 0),
  CONSTRAINT "repository_analysis_contributions_repository_full_name_check" CHECK (length(repository_full_name) > 0)
);
-- Create index "repository_analysis_contributions_occurred_at_idx" to table: "repository_analysis_contributions"
CREATE INDEX "repository_analysis_contributions_occurred_at_idx" ON "repository_analysis_contributions" ("occurred_at" DESC);
-- Create index "repository_analysis_contributions_user_id_occurred_at_idx" to table: "repository_analysis_contributions"
CREATE INDEX "repository_analysis_contributions_user_id_occurred_at_idx" ON "repository_analysis_contributions" ("user_id", "occurred_at" DESC);
