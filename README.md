# Gator: A Simple Blog Aggregator

This project, Gator, is really just a silly little thing, created for a guided project on [boot.dev](https://www.boot.dev/), which I highly suggested to aspiring programmers.

## Prerequisites

You'll need [Go](https://go.dev/) and [Postgres](https://www.postgresql.org/) installed on your system. Refer to their documentation for specifics on installation. You should do this all on a UNIX system (WSL, actual Linux, Mac). Results may be odd if using on straight up Windows.

## Setup

Besides initializing whatever you need for Postgres in general, you're also going to need to create a database for this program to run. Call it 'gator'.

Then, you must created a file in your home directory (whatever Go's `os.UserHomeDir()` function would retrieve) called `.gatorconfig.json` that contains a valid connection url to that database, structured as follows with likely url format:

```json
{
    "db_url":"postgres://[username]:[password]@localhost:5432/gator?sslmode=disable"
}
```

Where username and password are your postgres username and password.

## Installation

You can install Gator with the following command:

```bash
go install github.com/mrmiffmiff/gator@latest
```

## Usage

Run commands as `gator <command> [args...]` (or `go run . <command> [args...]` from source). Commands marked "requires login" need a current user set via `login` or `register` first.

| Command | Args | Description |
| --- | --- | --- |
| `register` | `<name>` | Creates a new user with the given name and sets them as the current user. |
| `login` | `<name>` | Sets the current user to an already-registered user. |
| `reset` | | Deletes all users from the database (and, by cascade, their feeds and follows). Use with caution. |
| `users` | | Lists all registered users, marking the current user. |
| `addfeed` | `<name> <url>` | Adds a new feed and automatically follows it as the current user. *(requires login)* |
| `feeds` | | Lists all feeds saved in the database along with their owners. |
| `follow` | `<url>` | Follows an existing feed by URL as the current user. *(requires login)* |
| `following` | | Lists the feeds the current user follows. *(requires login)* |
| `unfollow` | `<url>` | Unfollows the given feed. *(requires login)* |
| `agg` | `<time_between_reqs>` | Continuously scrapes the least-recently-fetched feed at the given interval (e.g. `1m`, `30s`) and stores new posts. Runs until interrupted. |
| `browse` | `[post_limit]` | Shows recent posts from feeds the current user follows, newest first, defaulting to 2 posts if no limit is given. *(requires login)* |
