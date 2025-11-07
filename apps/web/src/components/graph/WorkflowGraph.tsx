import { useMemo } from "react";
import {
  ReactFlow,
  Background,
  type Node,
  type Edge,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import type { GraphPayload } from "@/api/generated";
import { layoutGraph } from "./graphLayout";
import { GraphLegend } from "./legend";

export function nodeStatusStyles(failed?: boolean) {
  return failed ? { border: "2px solid var(--color-danger)" } : {};
}

export function WorkflowGraph({
  graph,
  onSelectNode,
}: {
  graph: GraphPayload;
  onSelectNode?: (nodeId: string) => void;
}) {
  const laid = useMemo(() => layoutGraph(graph), [graph]);
  const nodes: Node[] = laid.nodes.map((n, i) => ({
    id: n.id,
    position: { x: (i % 4) * 180, y: Math.floor(i / 4) * 100 },
    data: { label: n.topic ?? n.id },
    style: nodeStatusStyles(n.failed),
  }));
  const edges: Edge[] = laid.edges.map((e, i) => ({
    id: `e-${i}`,
    source: e.from,
    target: e.to,
    animated: graph.nodes.find((n) => n.id === e.from)?.failed,
  }));

  return (
    <div style={{ height: 360 }}>
      <GraphLegend />
      <ReactFlow
        nodes={nodes}
        edges={edges}
        fitView
        onNodeClick={(_, n) => onSelectNode?.(n.id)}
      >
        <Background />
      </ReactFlow>
    </div>
  );
}
