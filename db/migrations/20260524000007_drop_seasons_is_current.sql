ALTER TABLE seasons DROP COLUMN is_current;
DROP INDEX IF EXISTS seasons_is_current_idx;
