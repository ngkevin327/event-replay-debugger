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
      <div className="usage-meter__header">
        <span className="usage-meter__label">Daily event usage</span>
        <span className="usage-meter__value">{percent}%</span>
        <span className="badge">Starter plan</span>
      </div>
      <div className="usage-meter__track">
        <div
          className="usage-meter__fill"
          data-warn={warn}
          style={{ width: `${Math.min(100, percent)}%` }}
        />
      </div>
      {warn && (
        <p className="usage-meter__warn" role="status">
          Approaching plan limit — consider upgrading retention.
        </p>
      )}
    </div>
  );
}
