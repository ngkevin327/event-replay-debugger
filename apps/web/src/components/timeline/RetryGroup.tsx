export function RetryGroup({ label }: { label: string }) {
  return (
    <div className="retry-group" role="group" aria-label={`Retry chain ${label}`}>
      <span className="generation-label">{label}</span>
    </div>
  );
}
