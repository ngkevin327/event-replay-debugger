# Golden fixture: duplicate payout race

Fintech scenario where a payment retry succeeds after the original handler already credited the ledger, producing duplicate settlement on replay.

## Files

- `events.json` — captured Kafka timeline
- `expected_divergence.json` — replay divergence expectation

## Usage

Registered in `test/golden/run_determinism_test.go` as `duplicate_payout_race`.
