import type { Agent } from "@/api/generated";

function statusBadge(status?: string) {
  if (status === "healthy") return "Healthy";
  if (status === "offline") return "Offline";
  return "Unknown";
}

export function AgentHealthCard({ agent }: { agent: Agent }) {
  return (
    <article className="agent-card" data-status={agent.status ?? "unknown"}>
      <h3>{agent.hostname ?? agent.id}</h3>
      <span className="badge">{statusBadge(agent.status)}</span>
    </article>
  );
}
