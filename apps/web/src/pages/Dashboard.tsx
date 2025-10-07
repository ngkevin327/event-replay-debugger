import { AgentHealthCard } from "@/components/AgentHealthCard";
import { UsageMeter } from "@/components/UsageMeter";
import { useAgentHealth, useRecentIncidents } from "@/api/hooks";
import { Link } from "react-router-dom";

function RecentIncidentsTable({ projectId }: { projectId: string }) {
  const { data, isLoading } = useRecentIncidents(projectId);
  if (isLoading) return <p>Loading incidents…</p>;
  const rows = data?.incidents ?? [];
  if (!rows.length) return <p>No incidents yet.</p>;
  return (
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        {rows.map((inc) => (
          <tr key={inc.id}>
            <td>
              <Link to={`/incidents/${inc.id}`}>{inc.id.slice(0, 8)}</Link>
            </td>
            <td>{inc.status}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

export function Dashboard() {
  const projectId = localStorage.getItem("replay_project_id") ?? "";
  const agents = useAgentHealth();
  return (
    <div>
      <h1>Dashboard</h1>
      <UsageMeter projectId={projectId} />
      <section>
        <h2>Agents</h2>
        <div className="agent-grid">
          {(agents.data?.agents ?? []).map((a) => (
            <AgentHealthCard key={a.id} agent={a} />
          ))}
        </div>
      </section>
      <section>
        <h2>Recent incidents</h2>
        <RecentIncidentsTable projectId={projectId} />
      </section>
    </div>
  );
}

export default Dashboard;
