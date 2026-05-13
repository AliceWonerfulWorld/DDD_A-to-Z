-- Create "cp_accounts" table
CREATE TABLE "cp_accounts" ("user_id" text NOT NULL, "balance" bigint NOT NULL DEFAULT 0, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("user_id"), CONSTRAINT "cp_accounts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "cp_accounts_balance_check" CHECK (balance >= 0));

-- Create "cp_ledger" table
CREATE TABLE "cp_ledger" ("id" text NOT NULL, "user_id" text NOT NULL, "amount" bigint NOT NULL, "type" text NOT NULL, "reason" text NOT NULL, "source_type" text NOT NULL, "source_id" text NOT NULL, "balance_after" bigint NOT NULL, "created_at" timestamptz NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "cp_ledger_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "cp_ledger_amount_check" CHECK (amount <> 0), CONSTRAINT "cp_ledger_type_check" CHECK (type IN ('earn', 'spend', 'adjust')), CONSTRAINT "cp_ledger_reason_check" CHECK (length(reason) > 0), CONSTRAINT "cp_ledger_source_type_check" CHECK (length(source_type) > 0), CONSTRAINT "cp_ledger_source_id_check" CHECK (length(source_id) > 0), CONSTRAINT "cp_ledger_balance_after_check" CHECK (balance_after >= 0), CONSTRAINT "cp_ledger_type_amount_check" CHECK ((type = 'earn' AND amount > 0) OR (type = 'spend' AND amount < 0) OR type = 'adjust'));

-- Create index "cp_ledger_user_id_created_at_idx" to table: "cp_ledger"
CREATE INDEX "cp_ledger_user_id_created_at_idx" ON "cp_ledger" ("user_id", "created_at" DESC);

-- Create index "cp_ledger_source_idx" to table: "cp_ledger"
CREATE INDEX "cp_ledger_source_idx" ON "cp_ledger" ("source_type", "source_id");

-- Reject CP accounts that start with an implicit balance.
CREATE FUNCTION reject_nonzero_cp_account_insert()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  IF NEW.balance <> 0 THEN
    RAISE EXCEPTION 'cp account must start with zero balance'
      USING ERRCODE = '23514';
  END IF;

  RETURN NEW;
END;
$$;

CREATE TRIGGER cp_accounts_reject_nonzero_insert
BEFORE INSERT ON cp_accounts
FOR EACH ROW
EXECUTE FUNCTION reject_nonzero_cp_account_insert();

-- Only the CP ledger trigger may update account balances.
CREATE FUNCTION reject_direct_cp_account_balance_update()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  IF OLD.balance IS DISTINCT FROM NEW.balance
    AND pg_trigger_depth() < 2 THEN
    RAISE EXCEPTION 'cp account balance can only be updated from cp_ledger'
      USING ERRCODE = '23514';
  END IF;

  RETURN NEW;
END;
$$;

CREATE TRIGGER cp_accounts_reject_direct_balance_update
BEFORE UPDATE OF balance ON cp_accounts
FOR EACH ROW
EXECUTE FUNCTION reject_direct_cp_account_balance_update();

-- Create function and trigger that apply ledger entries to account balances.
CREATE FUNCTION apply_cp_ledger_entry()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
  next_balance BIGINT;
BEGIN
  SELECT balance + NEW.amount
  INTO next_balance
  FROM cp_accounts
  WHERE user_id = NEW.user_id
  FOR UPDATE;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'cp account not found for user_id %', NEW.user_id
      USING ERRCODE = '23503';
  END IF;

  IF next_balance < 0 THEN
    RAISE EXCEPTION 'cp balance cannot be negative for user_id %', NEW.user_id
      USING ERRCODE = '23514';
  END IF;

  UPDATE cp_accounts
  SET balance = next_balance,
      updated_at = NEW.created_at
  WHERE user_id = NEW.user_id;

  NEW.balance_after = next_balance;
  RETURN NEW;
END;
$$;

CREATE TRIGGER cp_ledger_apply_before_insert
BEFORE INSERT ON cp_ledger
FOR EACH ROW
EXECUTE FUNCTION apply_cp_ledger_entry();

-- Keep the CP ledger append-only after insertion.
CREATE FUNCTION reject_cp_ledger_mutation()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  RAISE EXCEPTION 'cp_ledger is append-only'
    USING ERRCODE = '23514';
END;
$$;

CREATE TRIGGER cp_ledger_reject_update
BEFORE UPDATE ON cp_ledger
FOR EACH ROW
EXECUTE FUNCTION reject_cp_ledger_mutation();

CREATE TRIGGER cp_ledger_reject_delete
BEFORE DELETE ON cp_ledger
FOR EACH ROW
EXECUTE FUNCTION reject_cp_ledger_mutation();
