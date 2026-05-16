-- Create "point_types" table
CREATE TABLE "point_types" ("code" text NOT NULL, "label" text NOT NULL, PRIMARY KEY ("code"));
-- Seed point types master data
INSERT INTO "point_types" ("code", "label") VALUES
  ('CP',             'Contribution Point'),
  ('Golang_SP',      'Golang Skill Point'),
  ('TypeScript_SP',  'TypeScript Skill Point');
-- Create "point_accounts" table
CREATE TABLE "point_accounts" ("user_id" text NOT NULL, "point_type" text NOT NULL, "balance" bigint NOT NULL DEFAULT 0, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("user_id", "point_type"), CONSTRAINT "point_accounts_point_type_fkey" FOREIGN KEY ("point_type") REFERENCES "point_types" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "point_accounts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "point_accounts_balance_check" CHECK (balance >= 0));
-- Create "point_ledger" table
CREATE TABLE "point_ledger" ("id" text NOT NULL, "user_id" text NOT NULL, "point_type" text NOT NULL, "amount" bigint NOT NULL, "type" text NOT NULL, "reason" text NOT NULL, "source_type" text NOT NULL, "source_id" text NOT NULL, "balance_after" bigint NOT NULL, "created_at" timestamptz NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "point_ledger_point_type_fkey" FOREIGN KEY ("point_type") REFERENCES "point_types" ("code") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "point_ledger_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "point_ledger_amount_check" CHECK (amount <> 0), CONSTRAINT "point_ledger_balance_after_check" CHECK (balance_after >= 0), CONSTRAINT "point_ledger_check" CHECK (((type = 'earn'::text) AND (amount > 0)) OR ((type = 'spend'::text) AND (amount < 0)) OR (type = 'adjust'::text)), CONSTRAINT "point_ledger_reason_check" CHECK (length(reason) > 0), CONSTRAINT "point_ledger_source_id_check" CHECK (length(source_id) > 0), CONSTRAINT "point_ledger_source_type_check" CHECK (length(source_type) > 0), CONSTRAINT "point_ledger_type_check" CHECK (type = ANY (ARRAY['earn'::text, 'spend'::text, 'adjust'::text])));
-- Create index "point_ledger_source_idx" to table: "point_ledger"
CREATE INDEX "point_ledger_source_idx" ON "point_ledger" ("source_type", "source_id");
-- Create index "point_ledger_user_id_created_at_idx" to table: "point_ledger"
CREATE INDEX "point_ledger_user_id_created_at_idx" ON "point_ledger" ("user_id", "point_type", "created_at" DESC);
-- Triggers: reject non-zero initial balance on point_accounts
CREATE FUNCTION reject_nonzero_point_account_insert()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
  IF NEW.balance <> 0 THEN
    RAISE EXCEPTION 'point account must start with zero balance'
      USING ERRCODE = '23514';
  END IF;
  RETURN NEW;
END;
$$;
CREATE TRIGGER point_accounts_reject_nonzero_insert
BEFORE INSERT ON point_accounts
FOR EACH ROW EXECUTE FUNCTION reject_nonzero_point_account_insert();
-- Triggers: only point_ledger trigger may update account balances directly
CREATE FUNCTION reject_direct_point_account_balance_update()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
  IF OLD.balance IS DISTINCT FROM NEW.balance
    AND pg_trigger_depth() < 2 THEN
    RAISE EXCEPTION 'point account balance can only be updated from point_ledger'
      USING ERRCODE = '23514';
  END IF;
  RETURN NEW;
END;
$$;
CREATE TRIGGER point_accounts_reject_direct_balance_update
BEFORE UPDATE OF balance ON point_accounts
FOR EACH ROW EXECUTE FUNCTION reject_direct_point_account_balance_update();
-- Triggers: apply ledger entry — compute balance_after and update point_accounts atomically
CREATE FUNCTION apply_point_ledger_entry()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
DECLARE
  next_balance BIGINT;
BEGIN
  SELECT balance + NEW.amount
  INTO next_balance
  FROM point_accounts
  WHERE user_id = NEW.user_id AND point_type = NEW.point_type
  FOR UPDATE;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'point account not found for user_id % point_type %', NEW.user_id, NEW.point_type
      USING ERRCODE = '23503';
  END IF;

  IF next_balance < 0 THEN
    RAISE EXCEPTION 'point balance cannot be negative for user_id % point_type %', NEW.user_id, NEW.point_type
      USING ERRCODE = '23514';
  END IF;

  UPDATE point_accounts
  SET balance    = next_balance,
      updated_at = NEW.created_at
  WHERE user_id = NEW.user_id AND point_type = NEW.point_type;

  NEW.balance_after = next_balance;
  RETURN NEW;
END;
$$;
CREATE TRIGGER point_ledger_apply_before_insert
BEFORE INSERT ON point_ledger
FOR EACH ROW EXECUTE FUNCTION apply_point_ledger_entry();
-- Triggers: keep point_ledger append-only
CREATE FUNCTION reject_point_ledger_mutation()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
  RAISE EXCEPTION 'point_ledger is append-only'
    USING ERRCODE = '23514';
END;
$$;
CREATE TRIGGER point_ledger_reject_update
BEFORE UPDATE ON point_ledger
FOR EACH ROW EXECUTE FUNCTION reject_point_ledger_mutation();
CREATE TRIGGER point_ledger_reject_delete
BEFORE DELETE ON point_ledger
FOR EACH ROW EXECUTE FUNCTION reject_point_ledger_mutation();
-- atlas:nolint DS102
-- Drop "contribution_point_accounts" table
DROP TABLE "contribution_point_accounts";
-- atlas:nolint DS102
-- Drop "contribution_point_ledger" table
DROP TABLE "contribution_point_ledger";
