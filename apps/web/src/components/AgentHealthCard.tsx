import type { Agent } from "@/api/generated";
import { StatusBadge } from "@/components/StatusBadge";

export function AgentHealthCard({ agent }: { agent: Agent }) {
  const label = agent.status ?? "unknown";
  return (
    <article className="agent-card" data-status={label}>
      <h3>{agent.hostname ?? agent.id}</h3>
      <StatusBadge status={label === "healthy" ? "healthy" : label === "offline" ? "offline" : "unknown"} />
    </article>
  );
}
