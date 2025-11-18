export function DivergenceDiff({
  expected,
  actual,
}: {
  expected: string;
  actual: string;
}) {
  return (
    <div className="divergence-diff">
      <div>
        <h3>Expected</h3>
        <pre>{expected}</pre>
      </div>
      <div>
        <h3>Actual</h3>
        <pre>{actual}</pre>
      </div>
    </div>
  );
}
