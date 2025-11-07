import type { GraphPayload } from "@/api/generated";

export function layoutGraph(graph: GraphPayload): GraphPayload {
  // MVP: pass-through; dagre would position nodes in production
  return graph;
}
