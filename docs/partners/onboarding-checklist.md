# Design partner onboarding checklist

Week-by-week tasks for the first three design partners on Replay MVP.

## Week 0 — Access

- [ ] Staging URL and org invite issued
- [ ] Partner engineer has login + API key
- [ ] Kafka non-production cluster identified for capture
- [ ] Topic allowlist agreed (max 5 topics)

## Week 1 — Agent install

- [ ] Helm install `replay-agent` 0.1.0 in partner staging namespace
- [ ] Heartbeat visible on dashboard within 15 minutes
- [ ] Sample traffic captured (>1k events)
- [ ] Review [quickstart](https://docs.replay.example/quickstart)

## Week 2 — First incident

- [ ] Create incident for known outage window
- [ ] Timeline reviewed with partner on call
- [ ] Graph view used to identify failed handler
- [ ] Export JSON shared with compliance (if needed)

## Week 3 — First replay

- [ ] Strict replay run completed
- [ ] Divergence report reviewed (expected or unexpected)
- [ ] Webhook `replay.completed` received and signature verified
- [ ] Feedback survey submitted ([feedback-template.md](./feedback-template.md))

## Week 4 — Sign-off

- [ ] Playbook dry run: duplicate payout scenario
- [ ] Open P0/P1 issues triaged with Replay team
- [ ] Decision: expand topics / move to production pilot

## Contacts

| Role | Channel |
|------|---------|
| Partner engineer | Slack `#partner-<name>` |
| Replay support | support@replay.example |
