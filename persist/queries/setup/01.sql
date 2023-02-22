BEGIN;

CREATE TABLE kv_props (
  key TEXT PRIMARY KEY,
  value
) WITHOUT ROWID;

INSERT INTO kv_props VALUES
  ('schema_version', 1),
  ('rate_limit_cooldown', '1 hour');

CREATE TABLE attributes (
  id TEXT NOT NULL,
  attr TEXT NOT NULL
);

CREATE UNIQUE INDEX attributes_id_attr_idx ON attributes (id, attr);

CREATE TABLE rate_limit_event (
  id TEXT PRIMARY KEY,
  timestamp TEXT NOT NULL
) WITHOUT ROWID;

COMMIT;
