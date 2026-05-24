-- atlas:nolint DS103
ALTER TABLE seasons DROP COLUMN IF EXISTS is_current;
DROP INDEX IF EXISTS seasons_is_current_idx;
