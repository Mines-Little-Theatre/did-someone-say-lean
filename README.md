# Lean Bot

A Discord bot that **LOVES LEAN**.

## Usage

Set the `LEAN_TOKEN` environment variable to the authorization token (should look like `<token string>`).

Set the `LEAN_DB` environment variable to the connection string for postgresql database. The application will create this database if it does not exist.

If you need a version of lean bot that uses SQLite - please use this commit `a52892837d1383bd5d730c35d97fe926124cdb21`.

Run with `go run .`!
