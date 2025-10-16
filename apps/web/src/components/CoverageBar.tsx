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
      <div className="fill" style={{ width: `${percent}%` }} />
      <span>{percent.toFixed(0)}% in retention window</span>
    </div>
  );
}
