import { useState } from "react";
import { Link } from "react-router-dom";
import { Pagination } from "@/components/Pagination";
import { CreateIncidentModal } from "@/components/CreateIncidentModal";
import { PageHeader } from "@/components/PageHeader";
import { useIncidents } from "@/api/hooks";
import type { IncidentStatus } from "@/api/generated";

export function StatusFilter({
  value,
  onChange,
}: {
  value: string;
  onChange: (v: string) => void;
}) {
  return (
    <select
      aria-label="Filter by status"
      value={value}
      onChange={(e) => onChange(e.target.value)}
    >
      <option value="">All statuses</option>
      <option value="collecting">Collecting</option>
      <option value="ready">Ready</option>
      <option value="failed">Failed</option>
    </select>
  );
}

function DataTable({ rows }: { rows: { id: string; status: IncidentStatus }[] }) {
  if (!rows.length) {
    return (
      <div className="empty-state">
        <p>No incidents match this filter.</p>
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
          {rows.map((r) => (
            <tr key={r.id}>
              <td>
                <Link to={`/incidents/${r.id}`}>{r.id.slice(0, 8)}…</Link>
              </td>
              <td>
                <span className={`badge status-${r.status}`}>{r.status}</span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export function IncidentsPage() {
  const projectId = localStorage.getItem("replay_project_id") ?? "";
  const [status, setStatus] = useState("");
  const [offset, setOffset] = useState(0);
  const [open, setOpen] = useState(false);
  const { data, isLoading } = useIncidents(projectId, status || undefined);

  return (
    <div>
      <PageHeader
        title="Incidents"
        description="Time-bounded collections of captured events for debugging and replay."
        actions={
          <button type="button" className="btn btn--primary" onClick={() => setOpen(true)}>
            Create incident
          </button>
        }
      />

      <div className="toolbar">
        <StatusFilter value={status} onChange={setStatus} />
      </div>

      {isLoading ? (
        <p className="empty-state">Loading incidents…</p>
      ) : (
        <DataTable rows={data?.incidents ?? []} />
      )}

      <Pagination offset={offset} limit={20} onChange={setOffset} />
      <CreateIncidentModal open={open} onClose={() => setOpen(false)} />
    </div>
  );
}

export default function Incidents() {
  return <IncidentsPage />;
}
