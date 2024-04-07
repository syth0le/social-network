-- CREATE INDEX user_full_name_idx ON user_table (first_name, second_name);
CREATE EXTENSION pg_trgm;
CREATE INDEX user_first_name_idx ON user_table USING gin (first_name, gin_trgm_ops);
CREATE INDEX user_second_name_idx ON user_table USING gin (second_name, gin_trgm_ops);
