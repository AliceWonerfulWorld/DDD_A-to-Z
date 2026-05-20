-- apply_point_ledger_entry トリガー関数を point_type_code + language カラムに対応させる
CREATE OR REPLACE FUNCTION apply_point_ledger_entry()
RETURNS TRIGGER LANGUAGE plpgsql AS $$
DECLARE
  next_balance BIGINT;
BEGIN
  SELECT balance + NEW.amount
  INTO next_balance
  FROM point_accounts
  WHERE user_id = NEW.user_id AND point_type_code = NEW.point_type_code AND language = NEW.language
  FOR UPDATE;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'point account not found for user_id % point_type_code % language %', NEW.user_id, NEW.point_type_code, NEW.language
      USING ERRCODE = '23503';
  END IF;

  IF next_balance < 0 THEN
    RAISE EXCEPTION 'point balance cannot be negative for user_id % point_type_code % language %', NEW.user_id, NEW.point_type_code, NEW.language
      USING ERRCODE = '23514';
  END IF;

  UPDATE point_accounts
  SET balance    = next_balance,
      updated_at = NEW.created_at
  WHERE user_id = NEW.user_id AND point_type_code = NEW.point_type_code AND language = NEW.language;

  NEW.balance_after = next_balance;
  RETURN NEW;
END;
$$;
