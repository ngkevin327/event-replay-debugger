# Agent data handling and egress

## Egress

The capture agent sends only configured topics to the Replay ingestion endpoint over HTTPS. Traffic uses `X-Replay-Key` authentication; keys are never logged.

## PII and redaction

Default redaction masks password, SSN, and card-number patterns in payloads before upload. Customers may extend rules via `REPLAY_REDACT_RULES`.

## Retention

- Disk buffer files are ephemeral under `REPLAY_BUFFER_DIR` and drained on successful flush.
- Long-term retention follows organization plan quotas in the control plane.

## Enterprise review

Security teams should verify topic allowlists, network policies to ingestion, and API key rotation procedures documented in the control plane API.
