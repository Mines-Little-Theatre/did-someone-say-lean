BEGIN;

UPDATE kv_props SET value = 2 WHERE key = 'schema_version';

CREATE TABLE gigglesnort (
  word TEXT PRIMARY KEY,
  message TEXT NOT NULL
) WITHOUT ROWID;

COMMIT;
