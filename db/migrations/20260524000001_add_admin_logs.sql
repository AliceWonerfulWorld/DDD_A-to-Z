-- Create "admin_logs" table
CREATE TABLE IF NOT EXISTS "admin_logs" ("id" bigserial NOT NULL, "action" text NOT NULL, "target_type" text NOT NULL, "target_id" text NOT NULL, "payload" jsonb NOT NULL, "created_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create index "admin_logs_created_at_idx" to table: "admin_logs"
CREATE INDEX IF NOT EXISTS "admin_logs_created_at_idx" ON "admin_logs" ("created_at" DESC);
