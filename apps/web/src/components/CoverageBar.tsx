export function CoverageBar({ percent }: { percent: number }) {
  return (
    <div
      className="coverage-bar"
      role="progressbar"
      aria-valuenow={percent}
      aria-valuemin={0}
      aria-valuemax={100}
      aria-label={`Coverage ${percent}%`}
    >
      <div className="usage-meter__track">
        <div
          className="usage-meter__fill"
          style={{ width: `${Math.min(100, percent)}%` }}
        />
      </div>
      <span className="text-muted">{percent.toFixed(0)}% in retention window</span>
    </div>
  );
}
