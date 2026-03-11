/** Consistent status pill across incidents, replays, and agents */
export function StatusBadge({ status }: { status: string }) {
  const normalized = status.toLowerCase().replace(/\s+/g, "-");
  return (
    <span className={`badge status-${normalized}`} data-status={normalized}>
      {status}
    </span>
  );
}
