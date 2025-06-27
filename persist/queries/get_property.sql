-- ?1: key
-- returns zero or one row (any?): value

SELECT value FROM kv_props WHERE key = $1;
