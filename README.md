# accessdesk

Internal tool for requesting and approving access to datasets. A Go backend
(built on [`foundry`](https://github.com/lattiq/foundry)) plus a React
frontend, backed by Postgres.

## Run it

```bash
docker compose up -d              # postgres
CONFIG_FILE=./config.yaml go run . # backend, http://localhost:8080

cd web
cp .npmrc.example .npmrc          # add a GitHub PAT with read:packages scope
npm install
npm run dev                       # frontend, http://localhost:5173
```

Seeded users (password for both: `password`):

| Email | Role |
|---|---|
| bhuvan@lattiq.com | requester |
| guna@lattiq.com | approver |

## How it works

- **Backend** (`main.go`, `internal/`): Gin router wired via foundry's
  `service` package. Postgres via GORM (`internal/store`), JWT auth
  (`internal/auth`), request/decision state machine (`internal/service`).
  Structured errors (`{code, message, details}`) on every failure path.
- **Frontend** (`web/`): Vite + React + TS, wiring `@lattiq/webtk` (auth
  store, API client), `@lattiq/auth` (login, token refresh), and
  `@lattiq/design-system` (all UI). Two pages: dataset catalog with a
  request form, and a requests table (approver sees everyone's, requester
  sees only their own; approve/reject behind a confirm dialog).
- **Tracing**: OpenTelemetry spans print to stdout on every request
  (`tracing.stdout_exporter: true` in `config.yaml`) — no collector needed
  for local dev.

Full build notes for the frontend (package auth, wiring decisions,
deviations from the spec) are in [`UI_GUIDE.md`](./UI_GUIDE.md).

## Caveats

- `config.yaml`'s `jwt.secret_key` and Postgres credentials are dev-only
  placeholders — never reuse them anywhere real.
- OpenObserve + the OTel collector (shipping traces off stdout) aren't wired
  up yet — stdout tracing only, for now.
- CORS is locked to `http://localhost:5173` (the Vite dev origin) because
  webtk's requests carry credentials, which browsers refuse to pair with a
  wildcard origin. Update it if the frontend runs anywhere else.
- `web/.npmrc` and `web/.env` are gitignored (they'd otherwise carry a
  personal access token / local URLs) — copy the `.example` files and fill
  them in yourself.
