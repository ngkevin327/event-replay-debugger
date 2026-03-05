import { AgentHealthCard } from "@/components/AgentHealthCard";
import { UsageMeter } from "@/components/UsageMeter";
import { PageHeader } from "@/components/PageHeader";
import { useAgentHealth, useRecentIncidents } from "@/api/hooks";
import { Link } from "react-router-dom";

function RecentIncidentsTable({ projectId }: { projectId: string }) {
  const { data, isLoading } = useRecentIncidents(projectId);
  if (isLoading) return <p className="empty-state">Loading incidents…</p>;
  const rows = data?.incidents ?? [];
  if (!rows.length) {
    return (
      <div className="empty-state">
        <p>No incidents yet.</p>
        <Link to="/incidents" className="btn btn--primary" style={{ marginTop: "1rem" }}>
          Create your first incident
        </Link>
      </div>
    );
  }
  return (
    <div className="data-table-wrap">
      <table className="data-table">
        <thead>
          <tr>
            <th>Incident</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((inc) => (
            <tr key={inc.id}>
              <td>
                <Link to={`/incidents/${inc.id}`}>{inc.id.slice(0, 8)}…</Link>
              </td>
              <td>
                <span className={`badge status-${inc.status}`}>{inc.status}</span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export function Dashboard() {
  const projectId = localStorage.getItem("replay_project_id") ?? "";
  const agents = useAgentHealth();
  const agentList = agents.data?.agents ?? [];

  return (
    <div>
      <PageHeader
        title="Dashboard"
        description="Monitor capture agents, usage, and recent incidents for your project."
        actions={
          <Link to="/incidents" className="btn btn--primary">
            New incident
          </Link>
        }
      />

      <UsageMeter projectId={projectId} />

      <section className="section">
        <h2 className="section__title">Capture agents</h2>
        {agentList.length === 0 ? (
          <div className="empty-state">
            <p>No agents connected.</p>
            <Link to="/agents/setup" className="btn btn--secondary" style={{ marginTop: "1rem" }}>
              Set up agent
            </Link>
          </div>
        ) : (
          <div className="card-grid">
            {agentList.map((a) => (
              <AgentHealthCard key={a.id} agent={a} />
            ))}
          </div>
        )}
      </section>

      <section className="section">
        <h2 className="section__title">Recent incidents</h2>
        <RecentIncidentsTable projectId={projectId} />
      </section>
    </div>
  );
}

export default Dashboard;
