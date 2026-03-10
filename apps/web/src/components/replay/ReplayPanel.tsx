import { useState } from "react";
import { useCreateReplay } from "@/api/hooks";
import { useNavigate } from "react-router-dom";

export function TimingModeSelect({
  value,
  onChange,
}: {
  value: string;
  onChange: (v: string) => void;
}) {
  return (
    <div className="form-field" style={{ marginBottom: 0 }}>
      <label htmlFor="timing-mode">Timing mode</label>
      <select
        id="timing-mode"
        aria-label="Timing mode"
        value={value}
        onChange={(e) => onChange(e.target.value)}
      >
        <option value="strict">Strict (deterministic)</option>
        <option value="compressed">Compressed</option>
      </select>
    </div>
  );
}

export function ReplayPanel({ incidentId }: { incidentId: string }) {
  const [mode, setMode] = useState("strict");
  const create = useCreateReplay(incidentId);
  const nav = useNavigate();

  async function start() {
    const run = await create.mutateAsync(mode);
    nav(`/replays/${run.id}`);
  }

  return (
    <section className="replay-panel animate-in">
      <h2>Start replay</h2>
      <p style={{ color: "var(--color-text-secondary)", fontSize: "var(--text-sm)", marginBottom: "1.25rem" }}>
        Re-run captured traffic in a sandbox and compare outcomes against the original timeline.
      </p>
      <TimingModeSelect value={mode} onChange={setMode} />
      <button
        type="button"
        className="btn btn--primary"
        style={{ marginTop: "1.25rem" }}
        onClick={start}
        disabled={create.isPending}
      >
        {create.isPending ? "Starting…" : "Run replay"}
      </button>
    </section>
  );
}
