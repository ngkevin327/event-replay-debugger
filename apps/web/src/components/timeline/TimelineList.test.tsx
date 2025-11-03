import { describe, expect, it } from "vitest";
import { groupRows } from "./TimelineList";
import type { TimelineEvent } from "./types";

describe("TimelineList", () => {
  it("renders 1k mock rows quickly", () => {
    const events: TimelineEvent[] = Array.from({ length: 1000 }, (_, i) => ({
      event_id: `e-${i}`,
      topic: "t",
      partition: 0,
      offset: i,
    }));
    const start = performance.now();
    const rows = groupRows(events);
    expect(rows.length).toBe(1000);
    expect(performance.now() - start).toBeLessThan(100);
  });
});
