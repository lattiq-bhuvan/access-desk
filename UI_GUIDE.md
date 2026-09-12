# Accessdesk Frontend (M6) — how it was built

This documents the M6 frontend: what was installed, how auth/API wiring
works, the folder layout, how to run it, and every deviation from a literal
reading of the track PDF along with the reason for it.

## Where it lives

`web/` at the repo root, next to the Go module (`main.go`, `internal/`,
`go.mod`). A Vite/React project has its own `package.json`, `node_modules`,
and TypeScript project — putting it in a subdirectory keeps `go build ./...`
and `npm run build` from ever touching each other's files.

## Getting the private packages

`@lattiq/design-system`, `@lattiq/webtk`, and `@lattiq/auth` are published to
GitHub Packages (`https://npm.pkg.github.com`), not the public npm registry,
and none of the three repos commit a `dist/` — only `src/`. Installing them
requires a **classic** GitHub PAT with the `read:packages` scope (a
fine-grained token cannot authenticate to GitHub's npm registry; `gh`'s own
CLI token doesn't carry `read:packages` by default either).

To reproduce this locally:

1. github.com/settings/tokens/new → classic token → scope `read:packages`
   only → generate.
2. `cp web/.npmrc.example web/.npmrc` and replace the placeholder with your
   token.
3. `cd web && npm install`.

`web/.npmrc` is gitignored — the token never enters git history. Only
`.npmrc.example` (with a placeholder) is committed.

Installed versions match the track PDF exactly: `@lattiq/design-system@0.7.2`,
`@lattiq/webtk@0.3.0`, `@lattiq/auth@0.2.1`. React is pinned to `18.3.1`
(not Vite's default 19) because that's the version `hub` — the production
frontend named as the closest reference in the PDF — and the design-system's
own devDependencies are built and tested against; using it removes a whole
class of "works upstream, subtly broken here" risk for zero cost.

`@tanstack/react-table` and `zod` are pinned to the exact versions
`@lattiq/design-system` depends on internally (`8.21.3` and `3.25.76`).
Both packages are regular (non-peer) dependencies of design-system, so if
our own `package.json` requested newer majors, npm would install two
incompatible copies side by side and TypeScript would see two unrelated
`ColumnDef`/`ZodType` types — exactly what happened on the first build
attempt. Pinning collapses them to one deduped copy.

## How the process went (so the reasoning is legible, not just the result)

1. Cloned `design-system`, `webtk`, `auth`, and — as reference — `hub` and
   `target` (read-only, outside this repo) to read their actual exported
   APIs instead of guessing from the PDF's prose. `hub` in particular is
   copied near verbatim in several places (`app.config.ts` shape,
   `tailwind.config.js`, `ProtectedLayout`/`PublicLayout` structure) because
   the PDF calls it out as "your closest model."
2. Read `webtk/src/auth/store.ts` and `types.ts` line by line — this is
   where the *real* contract lives, not just the type names. That's how the
   following requirements surfaced, none of which are stated in the PDF:
   - `login()` POSTs `{email, password, access}` — `access` is a mandatory
     "which portal is this" identifier, not something our backend needs to
     validate; Gin's default JSON binding silently ignores extra fields, so
     no backend change was needed to tolerate it.
   - `AuthConfig.buildApiUrl` is a full URL builder, not a base-URL string
     (the module doc-comment's `apiBaseUrl` example is stale/wrong).
   - `getMe()` populates `auth.user` from `GET /v1/users/me`, and the store
     gates `isAuthenticated: true` on that call succeeding — login isn't
     "done" until `/me` returns.
   - `User.roles` is `roles?: string[]` (plural, array) — not the backend's
     internal `role: string`. This one *did* require a backend change (below).
   - Every request uses `credentials: 'include'`, which rules out
     `Access-Control-Allow-Origin: *` on the backend the moment credentials
     are involved.
3. Scaffolded Vite + React + TS, installed the pinned packages, then hit a
   sequence of real build failures and fixed each at the root cause (see
   "Gotchas" below) rather than working around them — `npm run build`
   (`tsc -b && vite build`) is clean with zero errors and zero suppressed
   warnings.
4. Ran the actual stack end-to-end: Postgres → Go backend → curl login →
   `/v1/users/me` → `/v1/datasets`, then the Vite dev server, confirming the
   CORS preflight response and every JSON shape match what the frontend code
   expects, byte for byte.

## Deviations from a literal reading of the PDF, and why

- **Backend: `GET /v1/users/me` now returns `{id, email, name, roles: []}`
  instead of the model's native `{id, email, name, role: "..."}`.** This is
  the one non-frontend change. `webtk`'s `User` type — which is not
  something we authored and were told not to reimplement — reads
  `roles?: string[]`. Without this, the "approver sees Approve/Reject
  actions" requirement in §2.4 has no signal to key off on the frontend, so
  it's not optional. Implemented as a `dto.MeResponse` shape built in the
  handler (`internal/handler/auth.go`); the single-role column and every
  other part of the auth/authz model built in M3/M4 is untouched.
- **Backend: `config.yaml`'s CORS block** changed from `allowed_origins:
  ["*"]` to `["http://localhost:5173"]`, and `allowed_methods`/
  `allowed_headers` were made explicit (added `PATCH`/`PUT`/`DELETE`,
  `Content-Type`/`Authorization`). Required the moment the frontend's
  `credentials: 'include'` fetches are in the picture — browsers reject a
  wildcard origin on a credentialed request outright, so `"*"` isn't a
  looser version of this setting, it's a non-functional one here.
- **No AlertDialog/confirm-dialog component.** `@lattiq/design-system`
  exports a `Dialog` primitive but no dedicated confirm/alert-dialog
  component (despite `@radix-ui/react-alert-dialog` sitting in its
  `package.json`). The Approve/Reject confirm step
  (`src/components/ConfirmDialog.tsx`) composes `Dialog` + `Button` — using
  the library's own primitive, not a new one — the same way any consumer
  of this design system would have to.
- **`QuickForm` couldn't be used for the "Request access" dialog form.**
  Its prop type (`QuickFormProps`) doesn't actually include `showCard`, so
  passing it is a type error even though `FormWrapper` underneath supports
  it. Used `useFormValidation` (the hook `QuickForm` itself calls) +
  `FormWrapper` directly instead — same library code path, just called one
  layer down so `showCard={false}` is reachable (needed so the form doesn't
  render its own nested Card inside the Dialog's already-card-like surface).
- **No React Query.** `hub` uses `@tanstack/react-query`, but it isn't in
  the PDF's named library list and the accessdesk backend has no cache
  invalidation subtleties `useEffect` + local state can't handle at this
  scale. Added only what the PDF names.
- **Metric-card counts (stretch) computed client-side.** The backend has no
  aggregate-counts endpoint, so the three `MetricCard`s on the Requests page
  (approver view only) read `total` off three `page_size=1` list calls
  (one per status) rather than adding a new backend endpoint outside the
  five specified domain routes.
- **OTP is neither implemented nor reachable**, per the PDF's own
  instruction: `primaryMethod: 'password'` is set on `LoginPage`'s `config`
  and `allowToggle: false`, and the backend's `/auth/v1/login` never returns
  202 — so `LoginWithOTP`'s branch in `@lattiq/auth`'s `LoginPage` is simply
  never entered. `/auth/v1/logout` stays a `204` stub per the PDF; the
  `PATCH /v1/users/me/settings` stub already built in M4 satisfies webtk's
  `updateUserSettings()` (which expects exactly a `204`, no body).

## Gotchas hit during the build (fixed at the cause)

- **`baseUrl` in `tsconfig.app.json`** is deprecated in the TypeScript
  version Vite's scaffold pulled; `paths` resolves against the config file's
  own directory without it, so it was just dropped rather than suppressed.
- **`ColumnDef<AccessRequest>` alone doesn't compile** against
  `@tanstack/react-table@8.21.3`'s type definitions — it needs the second
  (`TValue`) type argument explicit: `ColumnDef<AccessRequest, unknown>`.
- **`@import '@lattiq/design-system/styles/globals.css'` failed the Tailwind
  build** with `the font-primary class does not exist` — the stylesheet's
  own `@apply font-primary` assumes the *consuming* app's `tailwind.config.js`
  defines a `fontFamily.primary` token (design-system ships Tailwind
  *source*, expecting to be compiled by the app that imports it, not a
  precompiled CSS bundle). Copied the exact `fontFamily`/`colors`/
  `borderRadius`/`keyframes` extensions `hub`'s config defines, since
  design-system's components assume that palette exists.
- **Two copies of `zod` and `@tanstack/react-table`** got installed (design-
  system's own pinned versions nested, plus whatever our `package.json`
  requested) — this is what produced the `showCard`-adjacent
  `QuickFormProps` mismatch and the `ColumnDef` arity error above. Pinning
  our versions to match design-system's collapsed both to one deduped copy
  each (verify with `npm ls zod @tanstack/react-table` — every line should
  say `deduped` or match).

## Folder layout

```
web/
  src/
    config/app.config.ts        platform name + webtk `access` identifier + storage key
    lib/api/
      baseUrl.ts                 VITE_API_BASE_URL -> full URL builder
      client.ts                  createApiClient() — the authenticated fetch wrapper
      errors.ts                  ApiError -> user-facing message / 401 check
      datasets.ts, requests.ts   thin wrappers over the 5 domain routes
    stores/useAuthStore.ts       createAuthStore() — zustand auth store
    providers/index.tsx          ThemeProvider + TokenRefreshProvider + Toaster
    components/
      layout/ProtectedLayout.tsx  AppLayout + role-aware nav + useLogout
      layout/PublicLayout.tsx     redirects an already-authenticated user away from /login
      ConfirmDialog.tsx           Dialog-based confirm step for Approve/Reject
    pages/
      Login.tsx                  wires @lattiq/auth's LoginPage, nothing else
      Catalog.tsx                Page 1 — dataset cards + request-access dialog/form
      Requests.tsx                Page 2 — data-table, status filter, approver actions, metric-cards
    types/domain.ts               Dataset / AccessRequest / DTO shapes matching the Go JSON tags
```

## Running it

```bash
# once, from repo root — Postgres + backend
docker compose up -d
CONFIG_FILE=./config.yaml go run .

# frontend
cd web
npm install   # after setting up web/.npmrc per "Getting the private packages" above
npm run dev   # http://localhost:5173
```

Log in as `bhuvan@lattiq.com` (requester) or `guna@lattiq.com` (approver),
password `password` for both (seeded in `internal/store/store.go`).

`npm run build` (`tsc -b && vite build`) passes clean — verified in this
session. The dev server was also started and its module graph exercised
(every page module fetched through Vite's transform pipeline with no
errors), and the full backend contract (`/auth/v1/login` →
`/v1/users/me` → `/v1/datasets`, plus a CORS preflight from the
`localhost:5173` origin) was verified against the running Go service with
curl, matching what the frontend code actually sends and expects byte for
byte.

**Not verified in this session: an actual browser render.** No browser was
available in this environment to click through the login form, submit a
request, or approve/reject one visually. Before your M6 checkpoint, open
`http://localhost:5173` yourself and walk both roles through: login →
Catalog → request access → toast → Requests page (as requester, own row
only) → log in as the approver → Requests page (all rows, Approve/Reject
behind the confirm dialog) → dark-mode toggle → log out. Also worth
forcing the 401/403/empty states the M6 checkpoint asks about directly:
stop the backend mid-session to see how the UI degrades, and check the
Catalog/Requests empty states by clearing the `access_requests` table.
