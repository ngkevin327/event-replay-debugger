export function GapMarker({
  gap,
}: {
  gap: { topic: string; partition: number; start_offset: number; end_offset: number };
}) {
  const label = `Gap ${gap.topic} p${gap.partition} offsets ${gap.start_offset}-${gap.end_offset}`;
  return (
    <div className="gap-marker" role="row" aria-label={label} title={label}>
      ⚠ Missing offsets {gap.start_offset}–{gap.end_offset}
    </div>
  );
}
