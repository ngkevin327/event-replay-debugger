import { Link } from "react-router-dom";
import type { ReplayRun } from "@/api/generated";
import { StatusBadge } from "@/components/StatusBadge";

export function ReplayHistory({ replays }: { replays: ReplayRun[] }) {
  const sorted = [...replays].sort((a, b) => (b.id > a.id ? 1 : -1));

  if (!sorted.length) {
    return (
      <div className="empty-state mt-4">
        <p>No replay runs yet. Start one above.</p>
      </div>
    );
  }

  return (
    <div className="data-table-wrap mt-4">
      <table className="data-table">
        <caption className="table-caption">Replay history</caption>
        <thead>
          <tr>
            <th>Run</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {sorted.map((r) => (
            <tr key={r.id}>
              <td>
                <Link to={`/replays/${r.id}`}>{r.id.slice(0, 8)}…</Link>
              </td>
              <td>
                <StatusBadge status={r.status} />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
