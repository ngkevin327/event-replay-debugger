# Playbook: Fintech duplicate payout race

Use this playbook when replay diverges on ledger credits after a payment retry storm.

## Symptoms

- Duplicate settlement detected in replay divergence report
- `ledger` topic shows two successful credits for one `payment_id`
- Golden fixture `test/golden/duplicate_payout_race` reproduces the issue

## Investigation steps

1. **Open the incident** in the web UI and confirm coverage > 90% for `payments` and `ledger` topics.
2. **Timeline view** — filter to `payments` and note retry generations (`gen1`, `gen2`).
3. **Graph view** — highlight failed nodes on the payment handler; follow cascade edges to ledger.
4. **Snapshots panel** — verify consumer offsets bracket the incident window.
5. **Start strict replay** — expect divergence at ledger offset ~50 (see fixture).

## Root cause patterns

| Pattern | Fix direction |
|---------|---------------|
| Idempotency key not propagated on retry | Add idempotent consumer keyed by `payment_id` |
| Race between timeout retry and late ack | Narrow timeout or use partition stickiness |
| Double publish from upstream | Dedupe at ingest with event_id |

## Share with compliance

Use **Export JSON** on the incident detail toolbar for audit attachments. Generate a **read-only share link** for external reviewers (72h TTL default).

## Related

- Golden fixture: `test/golden/duplicate_payout_race/`
- PRD fintech scenario — duplicate payout race
