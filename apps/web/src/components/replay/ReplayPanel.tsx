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
    <select
      aria-label="Timing mode"
      value={value}
      onChange={(e) => onChange(e.target.value)}
    >
      <option value="strict">Strict</option>
      <option value="compressed">Compressed</option>
    </select>
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
    <section className="replay-panel">
      <h2>Start replay</h2>
      <TimingModeSelect value={mode} onChange={setMode} />
      <button type="button" onClick={start} disabled={create.isPending}>
        Run replay
      </button>
    </section>
  );
}
