# Timeline performance (100k events)

Manual check for Stage 7 exit criteria:

1. Load incident detail with mock timeline of 100k events.
2. Scrub timeline range input — UI should stay responsive (16ms throttle).
3. Virtualized list via `react-window` keeps DOM node count low.

Target: p95 scrub latency under 500ms with warm cache.
