# FinTrack — Robust Personal Finance Tracker
**Requirements Document**
_Created: 2026-07-25_

---

## Goal

Transform FinTrack from a working prototype into a polished, daily-use personal finance tracker. This means fixing silent business logic bugs, overhauling the UI from the ground up, and adding the core missing interaction: editing transactions via modal. The result should feel like a premium personal finance tool the user actually wants to open every day.

## Primary Actor

Single user (self-hosted personal finance tracker). Auth already exists; no multi-tenancy changes required.

---

## Success Criteria

- User can edit any transaction without deleting and re-entering it
- Dashboard immediately communicates net worth and month-over-month change on load
- All business logic bugs identified in the codebase are fixed before any new features ship
- Searching and filtering transactions by date, category, type, and account takes under 2 seconds
- Transfer between accounts creates balanced ledger entries with no manual transaction pair
- The visual redesign passes the "looks like a real app" bar — consistent spacing, typography hierarchy, meaningful empty states

---

## Feature Areas

### 1. Business Logic Bug Fixes (ship first)

These are correctness issues, not features. Fix them before new UI lands so the redesign isn't built on broken data.

**a. Top Categories query returns 1 but renders as a list**
- `internal/managers/basic.manager.go` — `GetUserTopCategories` has `LIMIT 1` but the template iterates over a slice
- Fix: raise limit to 5, rename to reflect intent

**b. Bad error guard in home handler**
- `internal/handlers/gethome.go` lines 30 and 35: conditions are `err != nil && thisMonthExpenses > 0` — should be `err != nil`; current logic silently ignores real errors when the value happens to be positive

**c. Categories are global, not per-user**
- `internal/models/category.model.go` — `Category` has no `UserID`, so all users share the same category set
- Fix: add `UserID uuid.UUID` to `Category`; seed default categories per user on registration; migrate existing categories to the seeding user or treat them as system defaults
- `internal/seeders/categories.seeder.go` must be updated to accept a user ID

**d. CalculateBalance ignores Income when no prior transactions exist**
- `internal/managers/basic.manager.go` line 619-622: `if transactionType == constants.Income && len(transactions) != 0` — skips income on the very first transaction
- Fix: remove the `len(transactions) != 0` guard; income always adds to balance

**e. Hardcoded Indonesian strings**
- `CreateLoan` and `FinishLoan` in `internal/managers/basic.manager.go` use `"Memberi Hutang:"`, `"Hutang:"`, `"Bayar Hutang:"` as transaction descriptions
- Fix: use neutral English strings or make them configurable constants

**f. Recurring job closure captures loop variable**
- `internal/managers/basic.manager.go` `LoadAndScheduleJobs` — the `fn` closure captures `recurring` by reference in a loop, which is a Go closure bug (all jobs will fire with the last loop value)
- Fix: copy `recurring` into a local variable inside the loop before the closure

---

### 2. Transaction Edit Modal

**What:** Clicking any transaction row opens a modal with all editable fields. Saving updates the transaction and recalculates the affected account balance.

**Fields editable in modal:**
- Amount
- Description
- Category (dropdown, user's categories)
- Transaction date
- Account (dropdown — changing account recalculates both the old and new account balances)
- Transaction type (Income / Expense)

**Behavior:**
- Modal opens via HTMX (`hx-get="/transaction/{id}/edit"` → returns a templ partial)
- On save: `PUT /transaction/{id}` — manager updates record, recalculates balance(s)
- On success: HTMX swaps the transaction row in the list without full page reload
- If account changed: recalculate both old account and new account balances using `RecalculateAccountBalance`
- Cancel closes modal without changes

**New routes needed:**
- `GET /transaction/:id/edit` → returns edit modal partial
- `PUT /transaction/:id` → updates transaction, returns updated row partial

**Manager method needed:**
- `UpdateTransaction(id string, payload TransactionUpdateRequest) error` — handles balance recalculation for account changes

---

### 3. Transaction Search & Filter

**What:** The transactions page (`/transaction`) gains a persistent filter bar above the table.

**Filter dimensions:**
- Date range (start date / end date — two date inputs)
- Transaction type (All / Income / Expense — radio or segmented control)
- Category (dropdown, user's categories)
- Account (dropdown, user's accounts)
- Free-text search on description

**Behavior:**
- Filters are applied via HTMX on change (`hx-trigger="change"` on selects, `input delay:400ms` on text)
- Results replace the table body only, not the whole page
- Active filter count shown as a badge on the filter bar collapse button
- "Clear all" resets all filters

**New route:**
- `GET /transaction/list` — accepts query params for all filter dimensions, returns table body partial

**Pagination:**
- 20 transactions per page
- Pagination controls at bottom of table
- HTMX-driven page switching

---

### 4. Transfer Between Accounts

**What:** A "Transfer" action that moves money from one account to another, creating two balanced transactions.

**Behavior:**
- Accessible from the Accounts page and from a new "Transfer" button in the nav or quick-action area
- Form: From Account, To Account, Amount, Date, Description (optional)
- Creates two transactions atomically in a DB transaction:
  - Expense on the source account (category: "Transfer Out" — system category)
  - Income on the destination account (category: "Transfer In" — system category)
- Both transactions are linked by a shared `TransferGroupID uuid.UUID` field on the Transaction model so they can be identified and deleted together
- Deleting either transaction in the pair should delete both and reverse both balance changes

**New model field:** `TransferGroupID *uuid.UUID` on `Transaction` (nullable — nil for non-transfer transactions)

**New routes:**
- `GET /accounts/transfer` → transfer form partial
- `POST /accounts/transfer` → creates transfer pair, returns success

---

### 5. Dashboard Redesign — Net Worth First

**Layout direction:** Move from the current dense 5-column grid to a sidebar + main content layout.

**Navigation:** Left sidebar (fixed, collapsible) replacing the top horizontal nav links. Top nav becomes a minimal header with user info and a quick-add button.

**Dashboard sections (top to bottom, left to right):**

**Hero row:**
- Large "Net Worth" figure with month-over-month delta (e.g., `+Rp 2.4jt vs last month` in green)
- Month income / Month expense / Net — three smaller stat cards below

**Account cards row:**
- Horizontal scroll of account cards, each showing name, balance, account type icon
- "+ Add Account" card at the end

**Two-column lower section:**
- Left: Monthly income vs expense bar chart (6-month view, using Chart.js or ApexCharts loaded via CDN)
- Right: Top 5 spending categories (donut chart or ranked bar list)

**Recent transactions strip:**
- Last 5 transactions with edit/delete inline actions
- "View all" links to `/transaction`

**Upcoming recurring strip:**
- All upcoming recurring items in the next 30 days, not just the closest one
- Shows days remaining and amount

**Month-over-month delta calculation:**
- New manager method: `GetUserMonthlyTotals(userId string, year int, month int) (income, expense float64, error)`
- Called for current and previous month on home handler

---

### 6. Full UI Visual Language

**Typography:** Larger base size (16px body), stronger weight hierarchy — `font-extrabold` for key figures, `font-medium` for labels, `text-gray-400` for secondary info.

**Color palette:** Keep amber as the primary brand accent. Add a proper dark neutral surface for sidebar (`gray-900`). Cards: white with `shadow-sm` and `rounded-2xl`. Remove gradients on stat cards — replace with flat colored left-border accent or icon badge.

**Spacing:** Consistent 4/6/8 spacing scale. No more mixed gap-2 and gap-4 in the same level.

**Empty states:** Every list (accounts, transactions, loans, recurring) must have an illustrated or icon-based empty state with a clear CTA to add the first item.

**Loading feedback:** Add `htmx:beforeRequest` spinner overlay on forms. Disable submit buttons during in-flight requests.

**Toast / alert system:** Replace SweetAlert with a lightweight native toast (Templ partial injected into `#modals-here`) triggered by HX-Trigger headers. SweetAlert stays only for destructive confirm dialogs (delete, logout).

**Sidebar nav items:**
- Dashboard (Home icon)
- Transactions
- Recurring
- Loans
- Accounts
- [Divider]
- Profile / Settings (avatar + username)
- Logout

---

### 7. User-Specific Categories (from Bug Fix #c)

**What:** Categories belong to a user. On registration, the default seeded categories are cloned into the new user's set.

**Category management page (new, `/categories`):**
- List user's categories with edit name / delete
- Add custom category (name + optional color/icon selection)
- System categories ("Initial", "Transfer In", "Transfer Out") are shown as read-only

**Nav addition:** "Categories" link added to sidebar under Accounts.

---

## Scope Boundaries

### In scope
- All bug fixes listed in Feature Area 1
- Transaction edit modal + balance recalculation on account change
- Transaction search, filter, and pagination
- Transfer between accounts
- Dashboard visual redesign with sidebar, net worth hero, 6-month chart
- Per-user categories + category management page
- UI visual language overhaul across all pages
- Upcoming recurring shows all (not just closest)

### Deferred
- Budget limits per category (monthly spend cap with alert)
- Savings goals / targets
- CSV / PDF export
- Yearly periodicity for recurring (only Daily / Weekly / Monthly exist today)
- Mobile-optimized / responsive layout passes
- Push notifications for upcoming recurring
- User profile photo upload

### Outside this product's identity
- Multi-currency support
- Shared household / multi-user budgeting
- Bank API sync / import
- Mobile native app

---

## Key Decisions

| Decision | Resolution |
|---|---|
| Edit interaction pattern | Modal/drawer (HTMX partial), not inline or separate page |
| Dashboard anchor | Net worth hero + month-over-month delta |
| Navigation | Left sidebar replaces horizontal nav links |
| Category ownership | Per-user; system categories are seeded on registration |
| Transfer mechanism | Two linked transactions with `TransferGroupID` |
| Chart library | Chart.js or ApexCharts via CDN (no build step required) |
| Toast system | Native Templ partial replaces SweetAlert for non-destructive feedback |

---

## Assumptions

- The app remains single-user per deployment (no tenancy changes needed)
- Indonesian Rupiah (IDR) stays as the currency; formatting helpers in `internal/utils/currency.go` are kept
- HTMX + Templ stack stays — no React or SPA rewrite
- PostgreSQL schema changes are handled via GORM AutoMigrate (existing approach)
- Chart.js / ApexCharts loaded from CDN (consistent with current approach of CDN-loaded HTMX/Tailwind)
