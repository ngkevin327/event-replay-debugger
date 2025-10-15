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
    <nav aria-label="Pagination">
      <button
        type="button"
        disabled={offset === 0}
        onClick={() => onChange(Math.max(0, offset - limit))}
      >
        Previous
      </button>
      <span>
        Page {Math.floor(offset / limit) + 1}
      </span>
      <button type="button" onClick={() => onChange(offset + limit)}>
        Next
      </button>
    </nav>
  );
}
