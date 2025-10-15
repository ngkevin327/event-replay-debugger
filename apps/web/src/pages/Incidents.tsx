import { useState } from "react";
import { Link } from "react-router-dom";
import { Pagination } from "@/components/Pagination";
import { CreateIncidentModal } from "@/components/CreateIncidentModal";
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
      <option value="">All</option>
      <option value="collecting">Collecting</option>
      <option value="ready">Ready</option>
      <option value="failed">Failed</option>
    </select>
  );
}

function DataTable({ rows }: { rows: { id: string; status: IncidentStatus }[] }) {
  return (
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        {rows.map((r) => (
          <tr key={r.id}>
            <td>
              <Link to={`/incidents/${r.id}`}>{r.id.slice(0, 8)}</Link>
            </td>
            <td>{r.status}</td>
          </tr>
        ))}
      </tbody>
    </table>
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
      <h1>Incidents</h1>
      <StatusFilter value={status} onChange={setStatus} />
      <button type="button" onClick={() => setOpen(true)}>
        Create incident
      </button>
      {isLoading ? <p>Loading…</p> : <DataTable rows={data?.incidents ?? []} />}
      <Pagination offset={offset} limit={20} onChange={setOffset} />
      <CreateIncidentModal open={open} onClose={() => setOpen(false)} />
    </div>
  );
}

export default function Incidents() {
  return <IncidentsPage />;
}
