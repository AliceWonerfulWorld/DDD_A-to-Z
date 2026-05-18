CREATE TABLE point_types (
  code  TEXT PRIMARY KEY,
  label TEXT NOT NULL
);

INSERT INTO point_types (code, label) VALUES
  ('CP',        'Contribution Point'),
  ('Golang_SP', 'Golang Skill Point'),
  ('TypeScript_SP',     'TypeScript Skill Point');

CREATE TABLE point_accounts (
  user_id         TEXT NOT NULL REFERENCES users(id),
  point_type      TEXT NOT NULL REFERENCES point_types(code),
  balance         BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
  last_analyzed_at TIMESTAMPTZ,
  created_at      TIMESTAMPTZ NOT NULL,
  updated_at      TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (user_id, point_type)
);

CREATE TABLE point_ledger (
  id           TEXT PRIMARY KEY,
  user_id      TEXT NOT NULL REFERENCES users(id),
  point_type   TEXT NOT NULL REFERENCES point_types(code),
  amount       BIGINT NOT NULL CHECK (amount <> 0),
  type         TEXT NOT NULL CHECK (type IN ('earn', 'spend', 'adjust')),
  reason       TEXT NOT NULL CHECK (length(reason) > 0),
  source_type  TEXT NOT NULL CHECK (length(source_type) > 0),
  source_id    TEXT NOT NULL CHECK (length(source_id) > 0),
  balance_after BIGINT NOT NULL CHECK (balance_after >= 0),
  created_at   TIMESTAMPTZ NOT NULL,
  CHECK (
    (type = 'earn' AND amount > 0)
    OR (type = 'spend' AND amount < 0)
    OR type = 'adjust'
  )
);

CREATE INDEX point_ledger_user_id_created_at_idx ON point_ledger(user_id, point_type, created_at DESC);
CREATE INDEX point_ledger_source_idx ON point_ledger(source_type, source_id);

CREATE FUNCTION reject_nonzero_point_account_insert()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
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
FOR EACH ROW
EXECUTE FUNCTION reject_nonzero_point_account_insert();

CREATE FUNCTION reject_direct_point_account_balance_update()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
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
FOR EACH ROW
EXECUTE FUNCTION reject_direct_point_account_balance_update();

CREATE FUNCTION apply_point_ledger_entry()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
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
FOR EACH ROW
EXECUTE FUNCTION apply_point_ledger_entry();

CREATE FUNCTION reject_point_ledger_mutation()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  RAISE EXCEPTION 'point_ledger is append-only'
    USING ERRCODE = '23514';
END;
$$;

CREATE TRIGGER point_ledger_reject_update
BEFORE UPDATE ON point_ledger
FOR EACH ROW
EXECUTE FUNCTION reject_point_ledger_mutation();

CREATE TRIGGER point_ledger_reject_delete
BEFORE DELETE ON point_ledger
FOR EACH ROW
EXECUTE FUNCTION reject_point_ledger_mutation();
