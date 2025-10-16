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
    <div>
      <label>
        Start
        <input type="datetime-local" value={start} onChange={(e) => onStart(e.target.value)} />
      </label>
      <label>
        End
        <input type="datetime-local" value={end} onChange={(e) => onEnd(e.target.value)} />
      </label>
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
    <dialog open>
      <form onSubmit={submit}>
        <h2>Create incident</h2>
        {error && <p role="alert">{error}</p>}
        <DateRangePicker start={start} end={end} onStart={setStart} onEnd={setEnd} />
        <fieldset>
          <legend>Topics</legend>
          {(topics.data ?? []).map((t) => (
            <label key={t}>
              <input
                type="checkbox"
                checked={selected.includes(t)}
                onChange={(e) =>
                  setSelected((s) =>
                    e.target.checked ? [...s, t] : s.filter((x) => x !== t),
                  )
                }
              />
              {t}
            </label>
          ))}
        </fieldset>
        <button type="submit">Create</button>
        <button type="button" onClick={onClose}>
          Cancel
        </button>
      </form>
    </dialog>
  );
}
