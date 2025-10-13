import { useProjectUsage } from "@/api/hooks";

export function UsageMeter({ projectId }: { projectId: string }) {
  const { data } = useProjectUsage(projectId);
  const percent = data?.percent ?? 0;
  const warn = percent >= 80;
  return (
    <div
      className="usage-meter"
      role="meter"
      aria-valuenow={percent}
      aria-valuemin={0}
      aria-valuemax={100}
      aria-label="Daily event usage"
    >
      <span>Usage {percent}%</span>
      <div className="bar" data-warn={warn} style={{ width: `${percent}%` }} />
      {warn && <p role="status">Approaching plan limit</p>}
    </div>
  );
}
