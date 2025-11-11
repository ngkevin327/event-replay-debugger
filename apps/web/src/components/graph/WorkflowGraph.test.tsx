import { describe, expect, it } from "vitest";
import graphFixture from "./__fixtures__/graph.json";
import type { GraphPayload } from "@/api/generated";

describe("WorkflowGraph fixture", () => {
  it("mock graph JSON edge count", () => {
    const g = graphFixture as GraphPayload;
    expect(g.edges.length).toBe(3);
    expect(g.nodes.filter((n) => n.failed).length).toBeGreaterThan(0);
  });
});
