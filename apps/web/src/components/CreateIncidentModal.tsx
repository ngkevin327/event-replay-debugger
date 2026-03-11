import { FormEvent, useState } from "react";
import { useProjectTopics } from "@/hooks/useTopics";
import { apiFetch } from "@/api/client";
import { Modal } from "@/components/Modal";

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
    <div className="form-grid-2">
      <div className="form-field form-field--flush">
        <label htmlFor="inc-start">Window start</label>
        <input
          id="inc-start"
          type="datetime-local"
          value={start}
          onChange={(e) => onStart(e.target.value)}
        />
      </div>
      <div className="form-field form-field--flush">
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
    <Modal
      open={open}
      onClose={onClose}
      title="Create incident"
      footer={
        <>
          <button type="button" className="btn btn--ghost" onClick={onClose}>
            Cancel
          </button>
          <button type="submit" form="create-incident-form" className="btn btn--primary">
            Create incident
          </button>
        </>
      }
    >
      <form id="create-incident-form" onSubmit={submit}>
        {error && (
          <div className="alert alert--error" role="alert">
            {error}
          </div>
        )}
        <DateRangePicker start={start} end={end} onStart={setStart} onEnd={setEnd} />
        <fieldset>
          <legend>Topic filters</legend>
          <div className="checkbox-list">
            {(topics.data ?? ["payments.settlement"]).map((t) => (
              <label key={t} className="checkbox-item">
                <input
                  type="checkbox"
                  checked={selected.includes(t)}
                  onChange={(e) =>
                    setSelected((s) =>
                      e.target.checked ? [...s, t] : s.filter((x) => x !== t),
                    )
                  }
                />
                <span className="checkbox-item__label">{t}</span>
              </label>
            ))}
          </div>
        </fieldset>
      </form>
    </Modal>
  );
}
