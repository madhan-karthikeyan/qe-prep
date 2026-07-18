# Bug Reports

## Anatomy of a Good Bug Report

| Component | Why It Matters |
|-----------|---------------|
| **Title** | Clear, searchable summary of the problem |
| **Environment** | OS, browser, app version, deployment — context that affects reproducibility |
| **Steps to reproduce** | Exact sequence to trigger the bug. Minimal, numbered, complete. |
| **Expected result** | What should happen |
| **Actual result** | What actually happens |
| **Severity/Priority** | Helps triage and schedule |
| **Logs & Screenshots** | Evidence, not interpretation |
| **Additional context** | Frequency, workaround, related issues |

## Reproducible Steps

**Bad:**
> The login page crashes sometimes when I enter my email.

**Good:**
> 1. Go to https://app.example.com/login
> 2. Enter "alice@example.com" in the email field
> 3. Enter "TestPass123!" in the password field
> 4. Click "Sign In"
> 5. **Actual:** Page shows a white screen. Console error: `Uncaught TypeError: Cannot read properties of undefined (reading 'token')`
> 6. **Expected:** User is redirected to the dashboard

**Rule:** Someone who has never seen the app should be able to reproduce the bug by following steps 1–3.

## Environment Details

```
- OS: macOS 15.2 (Intel)
- Browser: Chrome 128.0.6613.84
- App Version: 2.5.1 (commit a3f2b1c)
- Database: PostgreSQL 16.3
- Feature flags: new_checkout=true
```

## Logs and Screenshots

- Attach raw log files, not screenshots of logs
- Annotate screenshots with arrows/boxes to highlight the problem
- For intermittent issues, include logs from both successful and failed runs
- Redact sensitive info (PII, tokens, secrets)

## Severity vs Priority

| | Critical | Major | Minor | Trivial |
|--|---------|-------|-------|---------|
| **High** | 🔥 Ship-blocking | Releasing without fix is risky | Should fix before next release | Fix if time permits |
| **Medium** | Must fix in this sprint | Fix in next sprint | Add to backlog | Consider future |
| **Low** | Depends on business impact | Defer | Deprioritize | Unlikely to fix |

**Severity** = how bad the impact is (data loss > cosmetic glitch)
**Priority** = how urgently it should be fixed (blocking release > nice to have)

## Common Bug Report Mistakes

| Mistake | Why It's Bad | Fix |
|---------|-------------|-----|
| Vague title ("Login doesn't work") | Hard to search, triage, or prioritize | "Login returns 500 for SSO users with empty name" |
| Missing steps | Developer can't reproduce | Write exact steps |
| Report from memory | Steps may be incorrect | Reproduce immediately before writing |
| Opinions instead of facts | "The UI is slow" → Developer can't act | "Click to response takes 8 seconds" |
| One bug per report (or many?) | One per bug — easier to track, fix, close | Separate issues, even if related |
| No version/environment | May not be reproducible on different setups | Always include versions |

## Bug Advocacy (How to Get Bugs Fixed)

1. **Make it impossible to ignore.** A bug with clear reproduction steps, severity rating, and business impact gets fixed.
2. **Link to user impact.** "This affects 15% of new signups, costing ~$500/day in lost conversions."
3. **Provide a failing test.** Nothing convinces a developer faster than a red test with `assert False` they can run.
4. **Classify correctly.** Don't mark everything P0/critical — you'll be ignored.
5. **Follow up.** If a bug sits in backlog for 2+ sprints, escalate with updated impact data.
6. **Be respectful.** Assume good faith. The developer may have constraints you don't see.

## Template

```markdown
## Bug: [Short descriptive title]

### Environment
- **App Version:** [version / commit hash]
- **OS:** [e.g., macOS 15.2, Ubuntu 24.04]
- **Browser:** [e.g., Chrome 128]
- **Database:** [e.g., PostgreSQL 16]

### Steps to Reproduce
1. [First step]
2. [Second step]
3. [Third step]
4. ...

### Expected Result
[What should happen]

### Actual Result
[What actually happens — include error messages]

### Severity / Priority
- **Severity:** [Critical / Major / Minor / Trivial]
- **Priority:** [High / Medium / Low]

### Logs
```
[paste relevant logs or attach files]
```

### Screenshots / Video
[attach or link]

### Additional Context
- **Frequency:** [Always / Intermittent — ~3 in 10]
- **Workaround:** [If any]
- **Related Issues:** [#123, #456]
```

### Example

```markdown
## Bug: Checkout fails with 500 when shipping address contains ampersand

### Environment
- **App Version:** 2.5.1 (a3f2b1c)
- **OS:** All
- **Browser:** All
- **Database:** PostgreSQL 16

### Steps to Reproduce
1. Add any item to cart
2. Go to checkout
3. Enter "123 Main St & 2nd Ave" in shipping address
4. Click "Place Order"

### Expected Result
Order is placed successfully, user sees confirmation page.

### Actual Result
Server returns 500. Error in logs: `ERROR: invalid input syntax for type json at character 123 — unescaped '&'`

### Severity / Priority
- **Severity:** Major
- **Priority:** High

### Logs
```
[2025-03-15 10:23:45] ERROR [checkout:42] — PG::InvalidJsonText: ERROR:  invalid input syntax
```

### Additional Context
- **Frequency:** Always
- **Workaround:** Remove & from address
- **Root cause:** Address JSON serialization doesn't escape special characters
```
```
