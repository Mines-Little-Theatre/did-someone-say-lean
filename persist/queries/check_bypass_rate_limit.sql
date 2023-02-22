-- ?1: user id
-- ?2: channel id
-- returns one row (integer): >0 if limits should be bypassed

SELECT EXISTS (
  SELECT * FROM attributes
  WHERE (id = ?1 OR id = ?2) AND attr = 'bypass_rate_limit'
);
