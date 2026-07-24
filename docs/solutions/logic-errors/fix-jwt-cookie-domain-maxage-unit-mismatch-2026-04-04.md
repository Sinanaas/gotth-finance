---
title: "Fix hardcoded cookie domain and MaxAge unit mismatch in JWT token refresh"
date: 2026-05-15
category: docs/solutions/logic-errors/
module: auth
problem_type: logic_error
component: authentication
symptoms:
  - After token refresh, session dropped within seconds on localhost
  - Refreshed cookies rejected by browser in staging, production, and Docker environments
  - Auth failures appear as missing cookie rather than invalid token
root_cause: config_error
resolution_type: code_fix
severity: high
tags:
  - jwt
  - cookies
  - gin
  - go
  - token-refresh
  - session-management
---

# Fix hardcoded cookie domain and MaxAge unit mismatch in JWT token refresh

## Problem

The `RefreshToken` method in `internal/managers/auth.manager.go` contained two bugs affecting cookie behavior after a token refresh: a hardcoded cookie domain that broke all non-localhost deployments, and a MaxAge unit mismatch that caused access tokens to expire 60× faster than intended.

## Symptoms

- After token refresh, session dropped within seconds on localhost
- Refreshed cookies rejected by browser in staging, production, and Docker environments
- Auth failures appear as "missing cookie" rather than "invalid token" — making the failure hard to distinguish from "never logged in"
- Login worked fine initially (the `Login` method had already received the domain and `*60` fixes), making the breakage appear intermittent and only reproducible after a token refresh event

## What Didn't Work

- **Checking token generation logic:** The JWT itself was generated and signed correctly — inspecting the token payload showed the right `exp` claim. This made it easy to conclude "the token is fine" and overlook that the *cookie's* `Max-Age` attribute was the problem.
- **Just fixing the domain:** Changing `"localhost"` to `""` would fix cross-environment deployments but leave the unit mismatch intact. Testers on localhost would still see sessions drop after ~15 seconds post-refresh.
- **Checking middleware/guard logic:** Because the cookie disappeared rather than containing an invalid token, auth failures looked like "no cookie present," pointing debugging toward browser cookie handling rather than `SetCookie` arguments in the refresh handler.
- **Reading the config value:** `AccessTokenMaxAge` reads as `15` — a reasonable integer with no signal that `15` means 15 seconds to the browser.

## Solution

**Bug 1 — Hardcoded domain in `RefreshToken`**

```go
// BEFORE
ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge,    "/", "localhost", false, true)
ctx.SetCookie("logged_in",    "true",      config.AccessTokenMaxAge,    "/", "localhost", false, false)

// AFTER
ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "",          false, true)
ctx.SetCookie("logged_in",    "true",      config.AccessTokenMaxAge*60, "/", "",          false, false)
```

The `domain` argument was changed from `"localhost"` to `""`. Gin passes an empty string to `http.Cookie`, which causes the browser to scope the cookie to the current request host — correct in every environment without any configuration change.

**Bug 2 — MaxAge unit mismatch**

`config.AccessTokenMaxAge` is stored in **minutes**, but Gin's `SetCookie` writes the value directly as the HTTP `Max-Age` attribute, which browsers interpret as **seconds**. Multiplying by `60` at the call site converts correctly:

```go
// BEFORE — 15-minute token expires in 15 seconds
config.AccessTokenMaxAge       // e.g. 15 → browser sees Max-Age=15

// AFTER — 15-minute token expires in 900 seconds (15 minutes)
config.AccessTokenMaxAge * 60  // e.g. 15 → browser sees Max-Age=900
```

Both fixes were applied together in commit `bb6197e`.

## Why This Works

**Empty domain string:** RFC 6265 and Go's `net/http` both specify that omitting the `Domain` attribute causes the browser to bind the cookie to the exact host that sent the response. Gin implements this by writing no `Domain=` attribute when the argument is `""`. Hardcoding any hostname breaks this the moment the app runs anywhere other than that exact host.

**Seconds conversion:** The HTTP `Max-Age` directive (RFC 6265 §4.1.2.2) is defined in seconds. Gin's `SetCookie` accepts an `int` and writes it verbatim as `Max-Age=<value>`. Because the config layer stores token lifetimes in minutes, every `SetCookie` call with a positive lifetime must multiply by `60`.

## Prevention

- **Active bug — `logged_in` cookie in `Login` is still missing `*60` (line 111 of `auth.manager.go`):**

  ```go
  // CURRENT — still broken in Login
  ctx.SetCookie("logged_in", "true", am.config.AccessTokenMaxAge, "/", "", false, false)

  // SHOULD BE
  ctx.SetCookie("logged_in", "true", am.config.AccessTokenMaxAge*60, "/", "", false, false)
  ```

  The `logged_in` indicator expires 60× faster than the actual `access_token` after a fresh login. Any UI or middleware reading `logged_in` will incorrectly show the user as logged out while a valid token is still present.

- **Define a seconds helper on Config** to eliminate unit confusion at every call site:

  ```go
  func (c *Config) AccessTokenMaxAgeSeconds() int {
      return c.AccessTokenMaxAge * 60
  }
  ```

- **Audit all `SetCookie` calls** for consistent `*60` application and `""` domain — the same class of mistake could recur anywhere cookies are set.

## Related Issues

- Commit `bb6197e` — "feat: styling and fixing jwt"
