export function Pagination({
  offset,
  limit,
  onChange,
}: {
  offset: number;
  limit: number;
  onChange: (next: number) => void;
}) {
  return (
    <nav className="pagination" aria-label="Pagination">
      <button
        type="button"
        className="btn btn--secondary"
        disabled={offset === 0}
        onClick={() => onChange(Math.max(0, offset - limit))}
      >
        Previous
      </button>
      <span className="pagination__info">Page {Math.floor(offset / limit) + 1}</span>
      <button type="button" className="btn btn--secondary" onClick={() => onChange(offset + limit)}>
        Next
      </button>
    </nav>
  );
}
