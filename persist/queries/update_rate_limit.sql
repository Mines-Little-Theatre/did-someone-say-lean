-- ?1: user id
-- ?2: channel id
-- returns no rows, sets the user's and channel's last rate limit events to now

INSERT INTO rate_limit_event VALUES
  (?1, datetime('now')),
  (?2, datetime('now'))
ON CONFLICT(id) DO UPDATE SET timestamp = excluded.timestamp;
