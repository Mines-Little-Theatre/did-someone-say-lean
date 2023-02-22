-- returns one row (integer): schema version

SELECT value FROM kv_props WHERE key = 'schema_version';
