# Database Structure

Data is persisted in an SQLite database pointed to by the `LEAN_DB` environment variable. The application will automatically create this database if the file does not exist.

## `kv_props` table

Store arbitrary values by text keys.

- **`schema_version`** (integer): The number of the last migration run against the database. The application can migrate a database with a lower version than it expects, but will refuse to run if the database has a higher version.
- **`rate_limit_cooldown`** (text): How long rate limits should last, stored as a [modifier](https://sqlite.org/lang_datefunc.html#modifiers). (default: `'1 hour'`)
- **`fallback_reaction`** (text): An emoji to react with if replying fails or is rate limited.
- **`gigglesnort_fallback_reaction`** (text): An additional emoji to use if the message would trigger gigglesnort. Maybe try that word again in the spam channel??

## `attributes` table

Associate textual attributes with users or channels.

- **`bypass_rate_limit`** (user, channel): Messages from this user/in this channel should ignore rate limit handling.
- **`ignore`** (user, channel): Messages from this user/in this channel should never be responded to.

## `rate_limit_event` table

Stores the timestamp of the last rate-limited event as returned by `datetime(now)`

## `gigglesnort` table

Associate a particular word (must be lowercase and include the substring "lean") with a particular message. The message MAY be cryptic as hell.
