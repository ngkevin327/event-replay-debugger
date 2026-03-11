import { DivergenceDiff } from "./DivergenceDiff";

export function MismatchRow({
  index,
  expected,
  actual,
}: {
  index: number;
  expected: string;
  actual: string;
}) {
  return (
    <tr data-mismatch={index === 0 ? "true" : undefined}>
      <td>{index}</td>
      <td>{expected}</td>
      <td>{actual}</td>
    </tr>
  );
}

export function DivergenceReport({
  mismatchIndex,
  expected,
  actual,
}: {
  mismatchIndex: number;
  expected: string;
  actual: string;
}) {
  return (
    <section className="divergence-report">
      <h2>Divergence at index {mismatchIndex}</h2>
      <div className="data-table-wrap">
        <table className="data-table">
          <thead>
            <tr>
              <th>Index</th>
              <th>Expected</th>
              <th>Actual</th>
            </tr>
          </thead>
          <tbody>
            <MismatchRow index={mismatchIndex} expected={expected} actual={actual} />
          </tbody>
        </table>
      </div>
      <DivergenceDiff expected={expected} actual={actual} />
    </section>
  );
}
