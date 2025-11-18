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
    <tr data-mismatch={index === 0}>
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
    <section>
      <h2>Divergence at index {mismatchIndex}</h2>
      <table>
        <tbody>
          <MismatchRow index={mismatchIndex} expected={expected} actual={actual} />
        </tbody>
      </table>
      <DivergenceDiff expected={expected} actual={actual} />
    </section>
  );
}
