export function SkeletonTimeline() {
  return (
    <div className="skeleton-timeline" aria-busy="true">
      <p className="text-muted">Collecting events…</p>
      <ul>
        {Array.from({ length: 5 }).map((_, i) => (
          <li key={i} className="skeleton-row" />
        ))}
      </ul>
    </div>
  );
}

export function ReconstructionProgress({ step }: { step: string }) {
  return (
    <p className="reconstruction-banner" aria-live="polite">
      Reconstruction in progress — <strong>{step}</strong>
    </p>
  );
}
