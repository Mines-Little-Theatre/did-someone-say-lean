BEGIN;

CREATE TABLE IF NOT EXISTS kv_props(
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL
);

INSERT INTO kv_props VALUES
  ('rate_limit_cooldown', '1 hour');

INSERT INTO kv_props VALUES ('user_version', '1');
INSERT INTO kv_props VALUES ('db_version', '1');

CREATE TABLE IF NOT EXISTS  attributes (
  id TEXT NOT NULL,
  attr TEXT NOT NULL
);

CREATE UNIQUE INDEX attributes_id_attr_idx ON attributes (id, attr);

CREATE TABLE IF NOT EXISTS rate_limit_event (
  id TEXT PRIMARY KEY,
  timestamp TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS gigglesnort (
  word TEXT PRIMARY KEY,
  message TEXT NOT NULL
);

COMMIT;
