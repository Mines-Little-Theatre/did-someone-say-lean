-- ?1: user id
-- ?2: channel id
-- returns one row (integer): >0 if a rate limit applies

WITH cooldown AS (SELECT value FROM kv_props WHERE key = 'rate_limit_cooldown')
SELECT EXISTS (
  SELECT * FROM rate_limit_event AS rl
  WHERE (rl.id = $1 OR rl.id = $2) AND datetime(rl.timestamp, cooldown.value) > datetime('now')
) FROM cooldown;
