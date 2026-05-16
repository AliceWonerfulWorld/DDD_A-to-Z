-- Create "contribution_point_accounts" table
CREATE TABLE "contribution_point_accounts" ("user_id" text NOT NULL, "balance" bigint NOT NULL DEFAULT 0, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("user_id"), CONSTRAINT "contribution_point_accounts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "contribution_point_accounts_balance_check" CHECK (balance >= 0));
-- Create "contribution_point_ledger" table
CREATE TABLE "contribution_point_ledger" ("id" text NOT NULL, "user_id" text NOT NULL, "amount" bigint NOT NULL, "type" text NOT NULL, "reason" text NOT NULL, "source_type" text NOT NULL, "source_id" text NOT NULL, "balance_after" bigint NOT NULL, "created_at" timestamptz NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "contribution_point_ledger_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "contribution_point_ledger_amount_check" CHECK (amount <> 0), CONSTRAINT "contribution_point_ledger_balance_after_check" CHECK (balance_after >= 0), CONSTRAINT "contribution_point_ledger_check" CHECK (((type = 'earn'::text) AND (amount > 0)) OR ((type = 'spend'::text) AND (amount < 0)) OR (type = 'adjust'::text)), CONSTRAINT "contribution_point_ledger_reason_check" CHECK (length(reason) > 0), CONSTRAINT "contribution_point_ledger_source_id_check" CHECK (length(source_id) > 0), CONSTRAINT "contribution_point_ledger_source_type_check" CHECK (length(source_type) > 0), CONSTRAINT "contribution_point_ledger_type_check" CHECK (type = ANY (ARRAY['earn'::text, 'spend'::text, 'adjust'::text])));
-- Create index "contribution_point_ledger_source_idx" to table: "contribution_point_ledger"
CREATE INDEX "contribution_point_ledger_source_idx" ON "contribution_point_ledger" ("source_type", "source_id");
-- Create index "contribution_point_ledger_user_id_created_at_idx" to table: "contribution_point_ledger"
CREATE INDEX "contribution_point_ledger_user_id_created_at_idx" ON "contribution_point_ledger" ("user_id", "created_at" DESC);
-- atlas:nolint DS102
-- Drop "cp_accounts" table (renamed to contribution_point_accounts; table is empty at migration time)
DROP TABLE "cp_accounts";
-- atlas:nolint DS102
-- Drop "cp_ledger" table (renamed to contribution_point_ledger; table is empty at migration time)
DROP TABLE "cp_ledger";
-- Drop orphaned cp_* functions (triggers were removed with their tables above)
DROP FUNCTION IF EXISTS reject_nonzero_cp_account_insert() CASCADE;
DROP FUNCTION IF EXISTS reject_direct_cp_account_balance_update() CASCADE;
DROP FUNCTION IF EXISTS apply_cp_ledger_entry() CASCADE;
DROP FUNCTION IF EXISTS reject_cp_ledger_mutation() CASCADE;
