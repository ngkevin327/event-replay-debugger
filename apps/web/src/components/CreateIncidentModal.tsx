import { FormEvent, useState } from "react";
import { useProjectTopics } from "@/hooks/useTopics";
import { apiFetch } from "@/api/client";

export function DateRangePicker({
  start,
  end,
  onStart,
  onEnd,
}: {
  start: string;
  end: string;
  onStart: (v: string) => void;
  onEnd: (v: string) => void;
}) {
  return (
    <div style={{ display: "grid", gap: "1rem", gridTemplateColumns: "1fr 1fr" }}>
      <div className="form-field" style={{ marginBottom: 0 }}>
        <label htmlFor="inc-start">Window start</label>
        <input
          id="inc-start"
          type="datetime-local"
          value={start}
          onChange={(e) => onStart(e.target.value)}
        />
      </div>
      <div className="form-field" style={{ marginBottom: 0 }}>
        <label htmlFor="inc-end">Window end</label>
        <input
          id="inc-end"
          type="datetime-local"
          value={end}
          onChange={(e) => onEnd(e.target.value)}
        />
      </div>
    </div>
  );
}

export function CreateIncidentModal({
  open,
  onClose,
}: {
  open: boolean;
  onClose: () => void;
}) {
  const projectId = localStorage.getItem("replay_project_id") ?? "";
  const topics = useProjectTopics(projectId);
  const [selected, setSelected] = useState<string[]>([]);
  const [start, setStart] = useState("");
  const [end, setEnd] = useState("");
  const [error, setError] = useState("");

  if (!open) return null;

  async function submit(e: FormEvent) {
    e.preventDefault();
    if (!selected.length) {
      setError("Select at least one topic");
      return;
    }
    await apiFetch(`/v1/projects/${projectId}/incidents`, {
      method: "POST",
      body: JSON.stringify({
        window_start: new Date(start).toISOString(),
        window_end: new Date(end).toISOString(),
        topic_filters: selected,
      }),
    });
    onClose();
  }

  return (
    <dialog open onClose={onClose}>
      <form className="modal-card" onSubmit={submit}>
        <h2>Create incident</h2>
        {error && (
          <div className="alert alert--error" role="alert">
            {error}
          </div>
        )}
        <DateRangePicker start={start} end={end} onStart={setStart} onEnd={setEnd} />
        <fieldset>
          <legend>Topic filters</legend>
          {(topics.data ?? ["payments.settlement"]).map((t) => (
            <label key={t} style={{ display: "flex", gap: "0.5rem", marginBottom: "0.5rem" }}>
              <input
                type="checkbox"
                checked={selected.includes(t)}
                onChange={(e) =>
                  setSelected((s) =>
                    e.target.checked ? [...s, t] : s.filter((x) => x !== t),
                  )
                }
              />
              <span style={{ fontFamily: "var(--font-mono)", fontSize: "var(--text-sm)" }}>
                {t}
              </span>
            </label>
          ))}
        </fieldset>
        <div style={{ display: "flex", gap: "0.75rem", justifyContent: "flex-end", marginTop: "1.5rem" }}>
          <button type="button" className="btn btn--ghost" onClick={onClose}>
            Cancel
          </button>
          <button type="submit" className="btn btn--primary">
            Create incident
          </button>
        </div>
      </form>
    </dialog>
  );
}
