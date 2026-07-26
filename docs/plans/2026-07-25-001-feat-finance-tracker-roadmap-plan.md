---
title: "feat: FinTrack roadmap — hardening, CRUD gaps, and core features"
type: feat
status: active
date: 2026-07-25
---

# feat: FinTrack roadmap — hardening, CRUD gaps, and core features

## Summary

Implements the full post-launch roadmap for FinTrack, a single-user, internet-hosted
personal finance tracker (Go + Gin + GORM + Postgres, Templ + HTMX + Tailwind). Work is
delivered in four phases: (1) go-live hardening — atomic loan writes, closed registration,
secure cookies, DB backups; (2) hardening + CRUD completeness — login rate limiting, IDOR
user-scoping, account/recurring editing, CSV export; (3) the features that make it a real
tracker — budgets, reports, savings goals, CSV import; (4) polish and durability — dark
mode, integration tests, operational hygiene. Phases are independently shippable and ordered
so data-integrity and hosting-safety land before net-new features.

---

## Problem Frame

FinTrack is functionally complete but was built for local, single-user use. The owner now
wants it **hosted on the internet for themselves only**. That shift changes what matters:
their own financial data must survive (backups), can't be silently corrupted (the loan
write bug), and can't be reached or trivially attacked by strangers who find the URL (open
registration, insecure cookies, no login rate limit). Separately, the app is still missing
the features that distinguish a "tracker" from a "ledger" — budgets, reports, goals — and
has two CRUD gaps (no account or recurring editing). This plan sequences all of it so the
hosting-safety and data-integrity work precedes feature growth.

Roadmap context established in prior brainstorm: audience = "just me, hosted online";
the key reframe was that **closing public registration** neutralizes most of the IDOR risk
far more cheaply than an exhaustive authorization audit, and that **CSV export doubles as a
manual backup / data-portability lever**.

---

## Requirements

- R1. Loan creation and completion write the loan and its transaction atomically; a partial failure leaves no orphaned records and no wrong balance.
- R2. Public registration can be disabled so only the owner has an account; when disabled, both the register form and the register endpoint refuse new signups.
- R3. Auth cookies are marked `Secure` when served over HTTPS, configurable so local HTTP dev still works.
- R4. The Postgres database is backed up automatically on a schedule, with timestamped dumps and bounded retention, restorable by a documented procedure.
- R5. Login (and registration, while enabled) is rate-limited to blunt brute-force attempts.
- R6. Every mutating operation is scoped to the logged-in user; a user cannot read, edit, or delete another user's record by supplying its id.
- R7. Accounts can be edited (name, description) after creation.
- R8. Recurring payments can be edited after creation, with their schedule updated to match.
- R9. Transactions can be exported to CSV.
- R10. Per-category monthly budgets can be set; each shows spend-to-date, remaining, and an over-budget state.
- R11. A reports view shows spending by category over a chosen range and net worth over time.
- R12. Savings goals can be created and tracked toward a target amount.
- R13. Transactions can be bulk-imported from a CSV file, with validation and a preview/error summary.
- R14. The UI supports a dark theme toggled by the user and persisted.
- R15. The manager/DB layer has integration test coverage for balance, transfer, and filter paths.
- R16. The server has structured logging, graceful shutdown, and fails fast on missing required config.

---

## Scope Boundaries

- Multi-user / multi-tenant hardening (per-tenant isolation guarantees, full CSRF token flow) — `SameSite=Lax` already shipped; single-user posture makes tokens unnecessary.
- Multi-currency and FX — remains IDR-only.
- Bank sync / open banking / statement OCR.
- Households, shared budgets, investment/portfolio tracking, debt-payoff planning beyond the existing loans feature.
- Native mobile app / PWA / push notifications (email/push infra out of scope; bill reminders deferred).

### Deferred to Follow-Up Work

- Bill/due-date reminders (needs email or push transport) — separate effort once a transport is chosen.
- Migrating the running container's deploy topology (reverse proxy / TLS provisioning) is an ops task tracked alongside U3/U4 but not code in this plan.

---

## Context & Research

### Relevant Code and Patterns

- Managers are split by domain under `internal/managers/` (`account.manager.go`, `transaction.manager.go`, `loan.manager.go`, `recurring.manager.go`, `category.manager.go`, `basic.manager.go`). All hang off `*BasicManager`.
- Balance is already unified: `RecalculateAccountBalance` + pure `SumBalance` in `internal/managers/account.manager.go` are the single source of truth. New write paths must recompute, never mutate incrementally.
- Handlers are one-file-per-route under `internal/handlers/`, wired in `internal/routers/basic.route.go` and `internal/routers/auth.route.go`. Mutations validate via `internal/validation` and respond with `swalError` / partial panels (`internal/handlers/helpers.go`).
- Templ components in `internal/templates/*.templ` (regenerate `_templ.go` via `templ generate`). Modal + partial-swap panel pattern already established (e.g., `AccountsPanel`, `addAccountModal`); reuse it for editing and budgets.
- Auth: `internal/managers/auth.manager.go` (`SignUp`, `Login`, `RefreshToken`, `Logout`), JWT RS256 + `gin-contrib/sessions` cookie store; the only middleware is `internal/middleware/deserialize_user.go`.
- Config via viper in `internal/initializers/loadEnv.go` (`Config` struct, `mapstructure` tags, `app.env`).
- Recurring scheduling uses `gocron/v2`: `SetUserCRONJob` / `RemoveCRONJob` / `LoadAndScheduleJobs` in `internal/managers/recurring.manager.go`.
- Currency helper: `internal/utils/currency.go` (`FormatCurrency`, IDR).
- Deploy: `Dockerfile`, `docker-compose.yml` (`app` + `db` services; `app.env` is untracked/gitignored). CI already present at `.github/workflows/go.yml` and `docker.yml`.
- Privacy/theme toggle pattern to mirror for dark mode: the `hide-balances` toggle in `internal/templates/layout.templ` (localStorage + pre-paint body class + top-bar button).

### Institutional Learnings

- `docs/solutions/logic-errors/` contains a prior write-up from the initial hardening pass (silent balance bugs). Same class of risk motivates R1/R6 — verify ownership and atomicity, don't trust incremental mutation.

### External References

- Go stdlib `encoding/csv` for export/import (R9/R13) — no dependency needed.
- Go `log/slog` (1.21+) for structured logging (R16).
- `testcontainers-go` (Postgres module) for hermetic DB integration tests (R15), since the schema relies on Postgres `uuid-ossp` and can't run on SQLite.
- Rate limiting can be done dependency-light with an in-memory fixed-window keyed by client IP; `didip/tollbooth` is an option if a maintained lib is preferred.

---

## Key Technical Decisions

- **Registration gating via config flag (`ALLOW_REGISTRATION`, default `false`).** Both `SignUp` and the `GET /register` route consult it; when off, the endpoint returns a refusal and the form isn't served. Simplest lockdown for a single-user host and reversible without code changes. Rationale: closing signup is the cheapest neutralizer of the IDOR attack chain (see origin reframe).
- **Cookie `Secure` from config (`COOKIE_SECURE`, default `false`).** Passed into every `ctx.SetCookie` and the session store options. Can't hardcode `true` (breaks local HTTP); can't hardcode `false` (unsafe when hosted). Env-driven is the standard resolution.
- **IDOR fix = thread `userId` into manager mutations and add `AND user_id = ?`.** Applies even though registration is closing — defense in depth, and cheap now that handlers already resolve `userId`. Delete/update return "not found" (not "forbidden") to avoid id-existence leakage.
- **Loan atomicity via `DB.Transaction(func(tx *gorm.DB) error { ... })`.** Replaces the cosmetic `Begin()`/`m.DB.Create` mix. All writes for one loan share `tx`; returning an error rolls back. Balance recompute runs after commit.
- **Backups via a compose sidecar running `pg_dump` on a schedule** to a mounted host volume with timestamped files and retention pruning. Keeps the backup lifecycle with the stack; documented restore path.
- **Budgets computed, not stored-as-spent.** A `Budget` row stores `(user, category, monthly limit)`; spend-to-date is derived from that category's expense transactions in the current month at read time — consistent with the "derive, don't mutate" balance philosophy.
- **Dark mode via Tailwind `darkMode: 'class'`** and a `dark` class toggled/persisted exactly like the existing `hide-balances` toggle. Adds `dark:` variants to the shared surfaces.
- **Reports reuse existing aggregates** (`GetUserMonthlyTotals`) plus one new category-breakdown query; net-worth-over-time is derived from month-end balances, no new stored history.

---

## Open Questions

### Resolved During Planning

- How to disable registration without ripping out auth? → Config flag gating both the endpoint and the form (U2); the owner's account already exists.
- Where do budgets get "spent" from? → Derived from transactions at read time, not a stored counter (U10).
- Can integration tests run in CI? → Yes, via testcontainers-go spinning a throwaway Postgres; gated behind a build tag so unit tests stay fast (U16).

### Deferred to Implementation

- Exact retention count / schedule for backups (e.g., daily, keep 14) — pick during U4 against the host's disk budget.
- Rate-limit thresholds (attempts/window) and whether to lock or just delay — tune in U5 against real login latency.
- Budget period beyond monthly (weekly/custom) — start monthly (R10); revisit only if needed.
- CSV column schema for import — finalize the header contract in U14 to match the export format from U9 (round-trip).

---

## High-Level Technical Design

> *This illustrates the intended approach and is directional guidance for review, not implementation specification. The implementing agent should treat it as context, not code to reproduce.*

Phase dependency shape (phases ship in order; units within a phase are mostly independent):

```
Phase 1 (go-live safety)        Phase 2 (hardening + CRUD)     Phase 3 (features)         Phase 4 (polish)
  U1 loan atomicity              U5 login rate limit            U10 budgets model          U15 dark mode
  U2 close registration  ─────►  U6 IDOR user-scoping   ─────►  U11 budgets UI      ─────►  U16 integration tests
  U3 secure cookies              U7 edit accounts               U12 reports                U17 ops hygiene
  U4 db backups                  U8 edit recurring              U13 savings goals
                                 U9 csv export ──────────────►  U14 csv import (round-trips U9)
```

New models introduced: `Budget` (U10), `Goal` (U13). New middleware: rate limiter (U5).
Everything else extends existing managers/handlers/templates.

---

## Implementation Units

### U1. Make loan create/finish atomic

**Goal:** Wrap `CreateLoan` and `FinishLoan` in a real DB transaction so a partial failure rolls back and never orphans a loan or writes a wrong `initial_transaction_id`.

**Requirements:** R1

**Dependencies:** None

**Files:**
- Modify: `internal/managers/loan.manager.go`
- Test: `internal/managers/loan_integration_test.go` (new; may land with U16's harness)

**Approach:**
- Replace the `tx := m.DB.Begin()` + `m.DB.Create(...)` mix with `m.DB.Transaction(func(tx *gorm.DB) error { ... })`, doing every write (loan insert, transaction insert, `initial_transaction_id` update) on `tx` and `return err` on any failure.
- Recompute the account balance with `RecalculateAccountBalance` **after** the transaction commits, not inside it.
- Remove the swallowed-error branches (`if err != nil { tx.Rollback() }` with no `return`).

**Patterns to follow:**
- `CreateTransfer` in `internal/managers/transaction.manager.go` (recompute-after-write); `SumBalance` as source of truth.

**Test scenarios:**
- Happy path: create a "lent" loan → one loan row, one expense transaction linked via `initial_transaction_id`, account balance decreased by amount.
- Happy path: finish a loan → reversing transaction created, loan soft-deleted/closed, balance restored.
- Error path: forced failure on the transaction insert → no loan row persists, balance unchanged (rollback proven).
- Integration: `initial_transaction_id` on the persisted loan equals the created transaction's id.

**Verification:** Injecting a failure mid-write leaves the DB unchanged; success produces exactly one loan + one linked transaction with a correct balance.

---

### U2. Close public registration behind a config flag

**Goal:** Allow the owner to disable new signups so only their account exists on the hosted instance.

**Requirements:** R2

**Dependencies:** None

**Files:**
- Modify: `internal/initializers/loadEnv.go` (add `AllowRegistration bool` / `ALLOW_REGISTRATION`)
- Modify: `internal/managers/auth.manager.go` (`SignUp` refuses when disabled)
- Modify: `internal/handlers/getregister.go` and `internal/handlers/postregister.go` (guard the endpoint)
- Modify: `internal/templates/login.templ` (hide/adjust the "create account" link when disabled) and `internal/templates/register.templ` (render a "registration closed" state)
- Modify: `app.env` (document the new var; keep untracked)
- Test: `internal/handlers/register_guard_test.go` (new)

**Approach:**
- Add the flag to `Config` (default `false` when absent).
- `POST /register` returns a refusal (swal error + 403) when the flag is off; `GET /register` renders a "registration is closed" message instead of the form.
- Thread the flag to handlers via the manager/config already available to auth handlers.

**Patterns to follow:**
- Config access pattern in `internal/initializers/loadEnv.go`; auth handler wiring in `internal/routers/auth.route.go`.

**Test scenarios:**
- Happy path (flag on): registration works exactly as today.
- Error path (flag off): `POST /register` with valid input → 403, no user created.
- Edge case (flag off): `GET /register` → renders closed-state, not the form.
- Edge case: flag absent from env → defaults to closed.

**Verification:** With `ALLOW_REGISTRATION=false`, no new account can be created by any request; with it `true`, current behavior is unchanged.

---

### U3. Secure cookies over HTTPS (config-driven)

**Goal:** Mark auth + session cookies `Secure` when hosted over TLS, without breaking local HTTP dev.

**Requirements:** R3

**Dependencies:** None

**Files:**
- Modify: `internal/initializers/loadEnv.go` (add `CookieSecure bool` / `COOKIE_SECURE`)
- Modify: `internal/managers/auth.manager.go` (pass secure flag to each `ctx.SetCookie`)
- Modify: `cmd/main.go` (session store `Options{Secure: ...}`)
- Modify: `app.env` (document var)

**Approach:**
- Replace the hardcoded `false` `secure` argument in every `ctx.SetCookie` call with the config value; set the session store's `Secure` option likewise.
- Keep `SameSite=Lax` (already shipped).

**Patterns to follow:**
- Existing `SetSameSite` + `SetCookie` calls in `internal/managers/auth.manager.go`; session `store.Options(...)` in `cmd/main.go`.

**Test scenarios:**
- Happy path (`COOKIE_SECURE=true`): login response sets `Secure` on access/refresh/session cookies.
- Edge case (`COOKIE_SECURE=false`): cookies omit `Secure` (local dev still works over HTTP).

**Verification:** Response `Set-Cookie` headers carry `Secure` iff the flag is set.

---

### U4. Automated Postgres backups

**Goal:** Scheduled, timestamped, retention-bounded DB dumps with a documented restore path.

**Requirements:** R4

**Dependencies:** None

**Files:**
- Modify: `docker-compose.yml` (add a `backup` sidecar; mount a `backups` volume)
- Create: `scripts/backup.sh` (pg_dump loop + retention prune)
- Create: `docs/operations/backup-restore.md` (procedure)

**Approach:**
- Sidecar (postgres-client image) runs `pg_dump` on a schedule against the `db` service, writing `gotth_finance-YYYYMMDD-HHMMSS.sql.gz` to a mounted host volume, pruning files older than the retention window.
- Restore documented as `gunzip | psql` against a fresh DB.
- Credentials come from `app.env` (already the compose env source).

**Patterns to follow:**
- Existing `db` service + `env_file: app.env` in `docker-compose.yml`.

**Test scenarios:**
- `Test expectation: none — ops/config; validated by running the sidecar and confirming a dump file appears and restores into a scratch DB (manual, documented in backup-restore.md).`

**Verification:** A dump file lands on schedule; restoring it into an empty DB reproduces the data; files beyond retention are pruned.

---

### U5. Rate-limit auth endpoints

**Goal:** Blunt brute-force on `POST /login` (and `POST /register` while enabled).

**Requirements:** R5

**Dependencies:** None

**Files:**
- Create: `internal/middleware/ratelimit.go`
- Modify: `internal/routers/auth.route.go` (apply to login/register)
- Test: `internal/middleware/ratelimit_test.go`

**Approach:**
- In-memory fixed-window limiter keyed by client IP (respecting proxy header only if a trusted-proxy setting is configured), returning 429 past the threshold. Single-instance host → no shared store needed.
- Thresholds configurable or sensible defaults; document in Deferred.

**Patterns to follow:**
- Gin middleware shape in `internal/middleware/deserialize_user.go`.

**Test scenarios:**
- Happy path: requests under the limit pass.
- Error path: N+1 requests within the window from one IP → 429.
- Edge case: window elapses → counter resets, requests allowed again.
- Edge case: different IPs tracked independently.

**Verification:** Exceeding the threshold from one client yields 429; normal login is unaffected.

---

### U6. Scope all mutations to the owning user (IDOR fix)

**Goal:** A user cannot edit/delete another user's account, transaction, loan, recurring, or category by supplying its id.

**Requirements:** R6

**Dependencies:** None

**Files:**
- Modify: `internal/managers/transaction.manager.go` (`DeleteTransactionById`, `UpdateTransaction`)
- Modify: `internal/managers/account.manager.go` (`DeleteAccountById`)
- Modify: `internal/managers/loan.manager.go` (`DeleteLoanById`, `FinishLoan`)
- Modify: `internal/managers/recurring.manager.go` (`DeleteRecurringById`)
- Modify: corresponding handlers in `internal/handlers/` to pass `userId`
- Test: `internal/managers/ownership_integration_test.go` (new; uses U16 harness)

**Approach:**
- Add a `userId` parameter to each mutation method and include `AND user_id = ?` in the `First`/`Delete`/`Save` queries. Handlers already resolve `userId` via `utils.GetSessionUserID`.
- Missing/foreign rows return a not-found error (surfaced as the existing swal error), not a distinct "forbidden".
- Category delete already scopes by `user_id` — confirm and align.

**Patterns to follow:**
- `DeleteUserCategory(categoryId, userId)` in `internal/managers/category.manager.go` (already scoped).

**Test scenarios:**
- Happy path: owner deletes/edits their own record → succeeds.
- Error path: user A supplies user B's transaction/account/loan/recurring id → not-found, B's data untouched.
- Edge case: nonexistent id → not-found, no panic.
- Integration: after a rejected cross-user delete, the target row still exists.

**Verification:** Every mutation query filters by `user_id`; a cross-user id cannot mutate data.

---

### U7. Edit accounts

**Goal:** Rename / edit an account's description after creation.

**Requirements:** R7

**Dependencies:** U6 (ownership scoping applies to the new update)

**Files:**
- Modify: `internal/managers/account.manager.go` (`UpdateAccount`)
- Create: `internal/handlers/geteditaccount.go`, `internal/handlers/patchaccount.go` (or `putaccount` edit variant)
- Modify: `internal/routers/basic.route.go` (routes)
- Create/Modify: `internal/templates/accounts.templ` (edit modal + row action)
- Test: `internal/managers/account_integration_test.go`

**Approach:**
- `UpdateAccount(userId, id, name, description)` scoped by user; does not touch balance (balance stays derived).
- Edit modal mirrors `addAccountModal`; the row gets a pencil action opening it pre-filled via a `GET /accounts/:id/edit` fragment; submit returns the refreshed `AccountsPanel`.

**Patterns to follow:**
- Edit-transaction modal flow (`GET /transaction/:id/edit` → `EditTransactionModal` → `PATCH`), partial-swap panels.

**Test scenarios:**
- Happy path: edit name/description → panel reflects change, balance unchanged.
- Error path: blank name → validation error (reuse `validation.Required`), no write.
- Error path (cross-user, via U6): editing another user's account id → not-found.

**Verification:** Editing updates only name/description for the owner's account and re-renders the panel.

---

### U8. Edit recurring payments

**Goal:** Edit a recurring payment and re-sync its schedule.

**Requirements:** R8

**Dependencies:** U6

**Files:**
- Modify: `internal/managers/recurring.manager.go` (`UpdateRecurring` + reschedule)
- Create: `internal/handlers/geteditrecurring.go`, `internal/handlers/patchrecurring.go`
- Modify: `internal/routers/basic.route.go`
- Modify: `internal/templates/recurring.templ` (edit modal + row action)
- Test: `internal/managers/recurring_integration_test.go`

**Approach:**
- `UpdateRecurring` updates the row (amount, type, category, account, name, periodicity, start date) scoped by user, then removes the existing gocron job and re-registers it (`RemoveCRONJob` + `SetUserCRONJob`) so the schedule matches.
- Reuse the recurring modal shape for the edit form.

**Patterns to follow:**
- `CreateRecurring` scheduling + `DeleteRecurringById`'s `RemoveCRONJob` in `internal/managers/recurring.manager.go`.

**Test scenarios:**
- Happy path: change amount/periodicity → row updated, old job removed, new job scheduled.
- Edge case: change of periodicity actually changes next occurrence (verify via `GetNextOccurrence`).
- Error path: invalid amount/date → validation error, no write, existing job intact.
- Error path (cross-user): not-found.

**Verification:** Editing updates the row and replaces the scheduled job so future runs use the new values.

---

### U9. CSV export of transactions

**Goal:** Download the owner's transactions as CSV (also serves as a manual backup).

**Requirements:** R9

**Dependencies:** U6 (export scoped to user)

**Files:**
- Create: `internal/handlers/getexporttransactions.go`
- Modify: `internal/routers/basic.route.go` (`GET /transactions/export`)
- Modify: `internal/templates/transaction.templ` (export button in the table header)
- Test: `internal/handlers/export_test.go`

**Approach:**
- Stream `text/csv` with `Content-Disposition: attachment`; columns: date, type, description, category, account, amount. Honor the same filters as the list when present (optional; default all).
- Use stdlib `encoding/csv`; scope query by `userId`.
- Define the header contract here so U14 import round-trips it.

**Patterns to follow:**
- `FilterTransactions` query in `internal/managers/transaction.manager.go` for the row set.

**Test scenarios:**
- Happy path: export → CSV with header row + one line per transaction, correct `Content-Type`/`Content-Disposition`.
- Edge case: no transactions → header-only CSV.
- Integration: exported rows belong only to the requesting user.

**Verification:** The download opens as a valid CSV containing exactly the owner's transactions.

---

### U10. Budget model + manager

**Goal:** Persist per-category monthly limits and compute spend-to-date.

**Requirements:** R10

**Dependencies:** U6 (ownership on budget mutations)

**Files:**
- Create: `internal/models/budget.model.go` (`Budget`: UserID, CategoryID, Amount, plus request structs)
- Create: `internal/managers/budget.manager.go` (`CreateBudget`, `UpdateBudget`, `DeleteBudget`, `GetUserBudgets`, `GetBudgetStatus`)
- Modify: `cmd/main.go` (AutoMigrate `&models.Budget{}`)
- Test: `internal/managers/budget_integration_test.go`

**Approach:**
- `Budget` = one monthly limit per (user, category). `GetBudgetStatus` returns, per budget, the limit, spend-to-date (sum of that category's **expense** transactions in the current month, derived — not stored), remaining, and an over-budget boolean.
- Reuse the month-window logic from `GetUserMonthlyTotals`.

**Patterns to follow:**
- `GetUserMonthlyTotals` month boundaries; category ownership rules from `category.manager.go`.

**Test scenarios:**
- Happy path: create budget, add expenses in category this month → status shows correct spent/remaining.
- Edge case: spend exceeds limit → over-budget true, remaining negative/zero per chosen convention.
- Edge case: transactions in prior months excluded from current spend.
- Error path: duplicate budget for same (user, category) → rejected or upserts (decide in unit; document).
- Integration: only the owner's transactions count toward their budget.

**Verification:** Budget status reflects real current-month category spend derived from transactions.

---

### U11. Budgets page + progress UI

**Goal:** A page to set/see budgets with progress bars and over-budget flags.

**Requirements:** R10

**Dependencies:** U10

**Files:**
- Create: `internal/handlers/getbudgets.go`, `internal/handlers/postbudget.go`, `internal/handlers/patchbudget.go`, `internal/handlers/deletebudget.go`
- Modify: `internal/routers/basic.route.go` (routes)
- Create: `internal/templates/budgets.templ` (`Budgets`, `BudgetsPanel`, `addBudgetModal`)
- Modify: `internal/templates/layout.templ` (sidebar nav link)
- Test: `internal/handlers/budget_handler_test.go`

**Approach:**
- Mirror the categories page: full-width panel + modal add/edit, partial-swap on mutation, validation via `internal/validation`.
- Each budget row shows category, limit (mono), spent, a progress bar (green→amber→red), and an over-budget badge.
- Respect the hide-balances blur (`font-mono`/`tnum`) and the theme.

**Patterns to follow:**
- `CategoriesPanel` + `addCategoryModal` in `internal/templates/categories.templ`; partial-swap handlers.

**Test scenarios:**
- Happy path: add a budget via modal → panel re-renders with the new row + progress.
- Edge case: over-budget row shows the red state + badge.
- Error path: missing category or non-positive amount → swal error, no write.
- Integration: progress reflects current-month spend (ties to U10).

**Verification:** The page lets the owner manage budgets and shows live progress per category.

---

### U12. Reports / analytics page

**Goal:** Spending-by-category over a range and net-worth-over-time.

**Requirements:** R11

**Dependencies:** None (reads existing data)

**Files:**
- Create: `internal/managers/report.manager.go` (`CategoryBreakdown(userId, start, end)`, `NetWorthSeries(userId, months)`)
- Create: `internal/handlers/getreports.go`
- Modify: `internal/routers/basic.route.go`
- Create: `internal/templates/reports.templ`
- Modify: `internal/templates/layout.templ` (nav link)
- Test: `internal/managers/report_integration_test.go`

**Approach:**
- Category breakdown = grouped sum of expenses per category over the range (reuse the top-categories query shape). Net-worth series = month-end balance across the last N months, derived from transactions.
- Render with Chart.js (already loaded): a category doughnut/bar + a net-worth line. Inject data via the existing `data-*` attribute pattern (no `{}` inside `<script>`).

**Patterns to follow:**
- `GetUserTopCategories` (grouped sum) and the dashboard chart data-attribute injection in `internal/templates/home.templ`.

**Test scenarios:**
- Happy path: breakdown sums expenses per category within the range; excludes income and out-of-range rows.
- Edge case: empty range → empty datasets, page renders without error.
- Integration: series values match month-end balances derived from the transaction set.

**Verification:** Reports render correct category breakdown and net-worth trend for the owner.

---

### U13. Savings goals

**Goal:** Create and track progress toward a savings target.

**Requirements:** R12

**Dependencies:** U6

**Files:**
- Create: `internal/models/goal.model.go` (`Goal`: UserID, Name, TargetAmount, optional linked AccountID, optional Deadline)
- Create: `internal/managers/goal.manager.go` (CRUD + progress)
- Create: handlers under `internal/handlers/` (`getgoals`, `postgoal`, `patchgoal`, `deletegoal`)
- Modify: `internal/routers/basic.route.go`, `cmd/main.go` (AutoMigrate), `internal/templates/layout.templ` (nav)
- Create: `internal/templates/goals.templ`
- Test: `internal/managers/goal_integration_test.go`

**Approach:**
- Progress = linked account balance toward target (if linked), else a manually updated current amount. Keep v1 simple: one number vs target + a progress bar.
- Same panel + modal + partial-swap pattern.

**Patterns to follow:**
- Budgets/categories page structure; `SumBalance` if linking to an account.

**Test scenarios:**
- Happy path: create goal → appears with 0→target progress; linked-account goal reflects that account's balance.
- Edge case: progress ≥ target → "reached" state.
- Error path: non-positive target / blank name → validation error.
- Integration (linked): editing the linked account's transactions moves goal progress.

**Verification:** Goals persist and show correct progress toward target.

---

### U14. CSV import of transactions

**Goal:** Bulk-create transactions from an uploaded CSV, with validation and an error/preview summary.

**Requirements:** R13

**Dependencies:** U9 (matching column contract), U6

**Files:**
- Create: `internal/handlers/postimporttransactions.go`
- Modify: `internal/routers/basic.route.go` (`POST /transactions/import`, multipart)
- Modify: `internal/managers/transaction.manager.go` (`ImportTransactions(userId, rows)` — atomic per file)
- Modify: `internal/templates/transaction.templ` (import modal: file input + result summary)
- Test: `internal/managers/import_integration_test.go`, `internal/handlers/import_handler_test.go`

**Approach:**
- Parse with `encoding/csv`; validate each row (amount>0, valid date, category/account resolvable **and owned by user**); collect per-row errors.
- Import valid rows in one `DB.Transaction`; return a summary (imported N, skipped M with reasons). Recompute affected account balances after commit.
- Round-trips U9's export format.

**Patterns to follow:**
- `internal/validation` helpers; `CreateTransaction` recompute-after-write; multipart via gin `ctx.FormFile`.

**Test scenarios:**
- Happy path: valid CSV → all rows imported, balances recomputed, summary counts correct.
- Edge case: mixed valid/invalid rows → valid imported, invalid reported with row numbers + reasons.
- Error path: category/account id not owned by user → row rejected (ties to U6).
- Error path: malformed CSV / missing header → 400 with a clear message, nothing imported.
- Integration: re-importing a file exported by U9 reproduces the same transactions.

**Verification:** A valid CSV imports atomically with a correct summary; bad rows are reported and skipped.

---

### U15. Dark mode

**Goal:** User-toggleable, persisted dark theme.

**Requirements:** R14

**Dependencies:** None

**Files:**
- Modify: `internal/templates/layout.templ` (Tailwind `darkMode: 'class'`, toggle button, persistence, pre-paint class)
- Modify: shared templates for `dark:` variants on surfaces (`home.templ`, panels, modals, `login.templ`, `register.templ`)
- Test: none (styling)

**Approach:**
- Set `darkMode: 'class'` in the inline Tailwind config; add a top-bar toggle mirroring the `hide-balances` toggle (localStorage `theme`, pre-paint `dark` class on `<body>`/`<html>` to avoid flash).
- Add `dark:` background/text/border variants to the cream/ink/stone surfaces; keep brass accent, tune for contrast.

**Patterns to follow:**
- The `hide-balances` toggle (localStorage + pre-paint body class + top-bar button) in `internal/templates/layout.templ`.

**Test scenarios:**
- `Test expectation: none — styling; verified visually and by confirming the class toggles and persists across reloads without a flash.`

**Verification:** Toggling switches theme instantly, persists across navigation/reload, no unstyled flash.

---

### U16. Integration test harness for the manager/DB layer

**Goal:** Hermetic Postgres-backed tests for balance, transfer, filter, and ownership paths.

**Requirements:** R15

**Dependencies:** None (unblocks/【validates U1, U6, U7, U8, U10, U13, U14 integration tests)

**Files:**
- Create: `internal/managers/main_test.go` (testcontainers Postgres setup, migrate schema, shared `*BasicManager`)
- Create/confirm: the `*_integration_test.go` files referenced by other units
- Modify: `.github/workflows/go.yml` (run integration tests, or gate behind a build tag + a CI job with a Postgres service)
- Test: (this unit *is* test infrastructure)

**Approach:**
- Spin a throwaway Postgres via testcontainers-go, run `AutoMigrate` + `uuid-ossp`, expose a manager bound to it. Guard behind a build tag (e.g., `//go:build integration`) so `go test ./...` stays fast by default; CI runs both.
- Seed helpers to create a user/account/category for scenarios.

**Patterns to follow:**
- Existing pure unit tests (`internal/validation/validation_test.go`, `internal/managers/balance_test.go`); `initializers.ConnectDB` migration list in `cmd/main.go`.

**Test scenarios:**
- Happy path: harness boots Postgres, migrates, and a trivial round-trip (create account → read back) passes.
- Integration: `RecalculateAccountBalance` over a seeded transaction set equals `SumBalance`.
- Integration: transfer between two accounts moves balances atomically.

**Verification:** `go test -tags=integration ./...` runs against a real Postgres and passes; default `go test ./...` stays unit-only and fast.

---

### U17. Operational hygiene: structured logging, graceful shutdown, config validation

**Goal:** Make the hosted server observable and fail-fast.

**Requirements:** R16

**Dependencies:** None

**Files:**
- Modify: `cmd/main.go` (`log/slog` logger; `http.Server` + `signal.NotifyContext` graceful shutdown; validate required config on boot)
- Modify: `internal/initializers/loadEnv.go` (return an error / fail fast on missing required keys)
- Test: `internal/initializers/config_validation_test.go`

**Approach:**
- Replace ad-hoc `log` with `slog` (JSON handler in prod). Wrap `server.Run` in an `http.Server` with graceful shutdown on SIGINT/SIGTERM so in-flight requests and the gocron scheduler stop cleanly.
- On boot, verify required config (DB creds, token keys, session secret) is present; exit with a clear message if not.

**Patterns to follow:**
- Current `cmd/main.go` init/boot flow; viper config load in `internal/initializers/loadEnv.go`.

**Test scenarios:**
- Happy path: all required config present → validation passes.
- Error path: missing session secret / token key → validation returns an error naming the missing key.
- Edge case: optional keys absent → no error.

**Verification:** Missing required config aborts boot with a clear message; shutdown drains cleanly; logs are structured.

---

## System-Wide Impact

- **Interaction graph:** New mutations (U7/U8/U11/U13/U14) go through the same validation → manager → partial-swap panel path; U6 changes every existing mutation's signature (handlers must pass `userId`). U8 and U14 touch the gocron scheduler and balance recompute respectively.
- **Error propagation:** Continue the established contract — mutation failures return `swalError` (422/4xx + `HX-Trigger` toast, no swap); GET failures render the error page; new auth refusals (U2/U5) use 403/429.
- **State lifecycle risks:** U1/U14 rely on `DB.Transaction` for atomicity; U8 must remove the old cron job before scheduling the new one (no duplicate jobs); budgets/goals are derived-at-read to avoid stored-counter drift.
- **API surface parity:** U6's `userId` scoping must be applied to *every* mutation, not a subset — audit the full handler list in `internal/routers/basic.route.go`.
- **Integration coverage:** U16 provides the harness the other units' integration tests depend on; land it early in Phase 4 or pull it forward if Phase 1–3 integration tests are wanted at implementation time.
- **Unchanged invariants:** Balance remains derived via `SumBalance`/`RecalculateAccountBalance` — no unit reintroduces incremental mutation. `SameSite=Lax` and the existing session/JWT scheme are unchanged except for the `Secure` flag (U3).

---

## Risks & Dependencies

| Risk | Mitigation |
|------|------------|
| U6 signature changes ripple across many handlers and could miss one, leaving an IDOR hole | Audit every mutation route in `internal/routers/basic.route.go`; add ownership integration tests (U16) covering each resource |
| Backup sidecar silently fails (no dumps) and is discovered only when a restore is needed | Document a periodic restore drill in `docs/operations/backup-restore.md`; log dump success/failure |
| Rate limiter keyed on a spoofable proxy header behind a reverse proxy | Only trust proxy headers when a trusted-proxy setting is configured; default to the direct remote addr |
| testcontainers unavailable in CI runner (no Docker-in-Docker) | Gate integration tests behind a build tag; run them in a CI job with a Postgres service container as fallback |
| Dark mode `dark:` variant sprawl across many templates → inconsistency | Centralize surface tokens; convert shared panels/modals first, verify visually before broad rollout |
| Loan atomicity refactor changes commit semantics and could alter `initial_transaction_id` timing | Cover with U1 integration tests asserting the link and rollback behavior |

---

## Documentation / Operational Notes

- New env vars (`ALLOW_REGISTRATION`, `COOKIE_SECURE`, rate-limit + backup settings) must be documented in `app.env` comments and a short "hosting" note.
- `docs/operations/backup-restore.md` is a deliverable (U4).
- Hosting prerequisites (TLS termination / reverse proxy so `COOKIE_SECURE=true` is meaningful) are an ops task adjacent to U3/U4.
- CI (`.github/workflows/go.yml`) gains an integration-test job (U16); ensure `templ generate` runs before build in CI if not already.

---

## Phased Delivery

### Phase 1 — Go-live hardening (ship before hosting)
- U1 (loan atomicity), U2 (close registration), U3 (secure cookies), U4 (backups). Protects existing data and closes the main hosted-attack surface. Independently shippable.

### Phase 2 — Hardening + CRUD completeness
- U5 (rate limit), U6 (IDOR scoping), U7 (edit accounts), U8 (edit recurring), U9 (CSV export). U6 is the backbone; U7/U8 depend on it.

### Phase 3 — Core features
- U10+U11 (budgets), U12 (reports), U13 (goals), U14 (CSV import; round-trips U9). The "real tracker" tier.

### Phase 4 — Polish + durability
- U15 (dark mode), U16 (integration harness — pull earlier if Phase 1–3 integration tests are wanted at build time), U17 (ops hygiene).

---

## Sources & References

- Prior brainstorm (roadmap): not persisted to a file; captured in this plan's Problem Frame and Key Technical Decisions.
- Related code: `internal/managers/`, `internal/handlers/`, `internal/templates/`, `internal/routers/basic.route.go`, `internal/initializers/loadEnv.go`, `cmd/main.go`.
- Institutional learnings: `docs/solutions/logic-errors/`.
- Prior completed work: `docs/brainstorms/robust-finance-tracker-requirements.md` (context only — that scope is shipped).
