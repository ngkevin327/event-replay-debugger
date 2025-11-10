export function GraphLegend() {
  return (
    <ul className="graph-legend" aria-label="Graph legend">
      <li>
        <span className="swatch healthy" /> Healthy node
      </li>
      <li>
        <span className="swatch failed" /> Failed node
      </li>
      <li>
        <span className="swatch cascade" /> Retry cascade edge
      </li>
    </ul>
  );
}
