# Database Structure

Data is persisted in an SQLite database pointed to by the `LEAN_DB` environment variable. The database needs to be manually created and populated before the application can be run.

## `kv_props` table

Store arbitrary values by text keys.

- **`schema_version`** (integer): The number of the last setup file run against the database. The application will not run if this is not exactly what it expects.
- **`rate_limit_cooldown`** (text): How long rate limits should last, stored as a [modifier](https://sqlite.org/lang_datefunc.html#modifiers). (default: `'1 hour'`)

## `attributes` table

Associate textual attributes with users or channels.

- **`bypass_rate_limit`** (user, channel): Messages from this user/in this channel should ignore rate limit handling.

## `rate_limit_event` table

Stores the timestamp of the last rate-limited event as returned by `datetime(now)`
