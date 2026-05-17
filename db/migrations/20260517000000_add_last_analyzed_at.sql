-- Add "last_analyzed_at" column to point_accounts for analysis idempotency
ALTER TABLE "point_accounts" ADD COLUMN "last_analyzed_at" timestamptz;
