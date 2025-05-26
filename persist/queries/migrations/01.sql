BEGIN;

CREATE TABLE kv_props (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL
);

INSERT INTO kv_props VALUES
  ('rate_limit_cooldown', '1 hour');

INSERT INTO kv_props VALUES ('user_version', '1');

CREATE TABLE attributes (
  id TEXT NOT NULL,
  attr TEXT NOT NULL
);

CREATE UNIQUE INDEX attributes_id_attr_idx ON attributes (id, attr);

CREATE TABLE rate_limit_event (
  id TEXT PRIMARY KEY,
  timestamp TEXT NOT NULL
);

CREATE TABLE gigglesnort (
  word TEXT PRIMARY KEY,
  message TEXT NOT NULL
);

COMMIT;
