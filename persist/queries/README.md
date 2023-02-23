# Database Structure

Data is persisted in an SQLite database pointed to by the `LEAN_DB` environment variable. The application will automatically create this database if the file does not exist.

## Metadata

The database should have an [application ID](https://www.sqlite.org/pragma.html#pragma_application_id) of `0x4c45414e` and a [user version](https://www.sqlite.org/pragma.html#pragma_user_version) corresponding to the last migration run against the database. The application will automatically migrate databases with too low a user version, but will refuse to run if the user version is too high.

## `kv_props` table

Store arbitrary values by text keys.

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

Wait. You're telling me there are SECRET `message`s that the bot might use in response to particular `word`s (store in all lowercase, must contain the substring "lean")? Holean shit. <!-- ;) -->
