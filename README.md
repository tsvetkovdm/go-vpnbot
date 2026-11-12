# go-vpnbot

A Go service that manages user subscriptions and access to an Xray proxy server.

- **Telegram bot** — users pick a plan and pay; admins confirm payments and manage users and plans.
- **HTTP API** — serves each user's client configuration by a personal subscription link.
- **Reconciler** — once a minute compares users in PostgreSQL with users in Xray and applies the difference via the Xray gRPC API.

## Stack

Go · PostgreSQL (pgx) · Xray gRPC API · Docker Compose · golang-migrate · zap

## How it works


- The bot, the HTTP server and the reconciler run in one process under `errgroup`: if one fails, the others shut down gracefully.
- Code is split by feature (`internal/features/...`), each with `transport → service → repository` layers. Features don't import each other; they talk through small interfaces wired in `main.go`.
- Concurrent updates are handled in the database: state transitions are conditions in `WHERE`, other changes use optimistic locking by `version`.
- The reconciler only touches users it created and supports a dry-run mode.

## Run

```bash
cp .env.example .env    # fill in the bot token and DB settings
make env-up             # start PostgreSQL
make migrate-up         # apply migrations
make vpnbot-deploy      # build and start the bot
```

Xray must be running with its API enabled; set its address in `XRAY_API_ADDR`.