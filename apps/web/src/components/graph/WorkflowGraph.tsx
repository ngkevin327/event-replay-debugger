import { useMemo, type CSSProperties } from "react";
import {
  ReactFlow,
  Background,
  Controls,
  type Node,
  type Edge,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import type { GraphPayload } from "@/api/generated";
import { useTheme } from "@/context/ThemeContext";
import { layoutGraph } from "./graphLayout";
import { GraphLegend } from "./legend";

export function nodeStyle(failed?: boolean): CSSProperties {
  return failed
    ? {
        background: "var(--graph-node-failed)",
        borderColor: "var(--color-danger)",
      }
    : {
        background: "var(--graph-node-bg)",
        borderColor: "var(--graph-node-border)",
      };
}

export function WorkflowGraph({
  graph,
  onSelectNode,
}: {
  graph: GraphPayload;
  onSelectNode?: (nodeId: string) => void;
}) {
  const { theme } = useTheme();
  const laid = useMemo(() => layoutGraph(graph), [graph]);
  const nodes: Node[] = laid.nodes.map((n, i) => ({
    id: n.id,
    position: { x: (i % 4) * 200, y: Math.floor(i / 4) * 110 },
    data: { label: n.topic ?? n.kind ?? n.id },
    style: nodeStyle(n.failed),
  }));
  const edges: Edge[] = laid.edges.map((e, i) => ({
    id: `e-${i}`,
    source: e.from,
    target: e.to,
    animated: graph.nodes.find((n) => n.id === e.from)?.failed,
  }));

  return (
    <div className="workflow-graph">
      <GraphLegend />
      <ReactFlow
        nodes={nodes}
        edges={edges}
        fitView
        colorMode={theme}
        onNodeClick={(_, n) => onSelectNode?.(n.id)}
        proOptions={{ hideAttribution: true }}
      >
        <Background gap={20} size={1} />
        <Controls showInteractive={false} />
      </ReactFlow>
    </div>
  );
}
