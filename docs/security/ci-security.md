# CI security gates

The `.github/workflows/security.yml` workflow runs on every pull request and weekly on `main`.

## Jobs

| Job | Tool | Fail threshold |
|-----|------|----------------|
| `dependency-scan` | Trivy filesystem scan | CRITICAL findings |
| `owasp-zap` | ZAP baseline against staging | High confidence alerts |

## Trivy

- Scans repository root including `apps/`, `services/`, and container Dockerfiles.
- Uploads SARIF to GitHub Security tab when `GITHUB_TOKEN` allows.
- **Blocks merge** on CRITICAL CVE with fix available.

## OWASP ZAP baseline

- Target: `STAGING_URL` repository variable (default `https://staging.replay.example`).
- Reports passive findings only; no active attack mode in MVP.
- **Blocks merge** on High risk passive findings in auth and session endpoints.

## Exemptions

Document exemptions in the PR description with ticket link and expiry date. Security oncall approves exceptions > 7 days.
