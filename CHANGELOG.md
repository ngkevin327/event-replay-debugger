# Changelog

All notable changes to the Replay Platform are documented in this file.

## [0.1.0-mvp] — 2026-05-20

### Added

- Kafka capture agent (Go) with Helm chart `replay-agent` 0.1.0
- Control plane API: auth, projects, API keys, incidents, timeline, graph, replay
- Ingestion pipeline with batch ingest and ClickHouse storage
- Reconstruction worker and timeline/graph artifacts
- Deterministic replay worker with divergence detection
- Web operator console (React): incidents, timeline, graph, replay UI
- Notifications service: signed webhooks and email templates
- Incident export and expiring read-only share links
- Terraform staging stack (VPC, EKS, RDS, S3, ClickHouse)
- Helm umbrella chart `replay-platform` and CI deploy workflows
- Security: rate limiting, tenant isolation tests, OWASP ZAP CI
- Public docs site (Docusaurus) with quickstart and fintech playbook

### Known limits (MVP)

- Single-region staging only; no multi-tenant data residency
- Manual billing / starter plan flags; no Stripe integration
- Timeline/graph loaders use stubs until artifacts wired in all environments
- WCAG audit scaffolded; full axe-core gate pending

See [docs/release/v0.1.0-mvp.md](docs/release/v0.1.0-mvp.md) for acceptance criteria mapping to PRD §19.
