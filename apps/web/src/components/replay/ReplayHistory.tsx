import { Link } from "react-router-dom";
import type { ReplayRun } from "@/api/generated";

export function ReplayHistory({ replays }: { replays: ReplayRun[] }) {
  const sorted = [...replays].sort(
    (a, b) => (b.id > a.id ? 1 : -1),
  );
  return (
    <table>
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
              <Link to={`/replays/${r.id}`}>{r.id.slice(0, 8)}</Link>
            </td>
            <td>{r.status}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
