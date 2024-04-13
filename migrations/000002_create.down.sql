DROP INDEX IF EXISTS user_first_name_idx;
DROP INDEX IF EXISTS user_second_name_idx;

DROP EXTENSION pg_trgm;
DROP EXTENSION btree_gin;