# E2E tests (Playwright)

## Prerequisites

- API gateway on `http://localhost:8080` (or mock routes in specs)
- `pnpm install` at repo root

## Run

```bash
pnpm exec playwright test
```

## Environment

- `CI` — enables retries and strict `forbidOnly`
- Web preview served on port 5173 via `playwright.config.ts` `webServer`
