import { useState } from "react";

export function JsonViewer({ data }: { data: unknown }) {
  const [open, setOpen] = useState(true);
  return (
    <div className="json-viewer">
      <button type="button" className="btn btn--ghost" onClick={() => setOpen(!open)}>
        {open ? "Collapse" : "Expand"} JSON
      </button>
      {open && <pre>{JSON.stringify(data, null, 2)}</pre>}
    </div>
  );
}
