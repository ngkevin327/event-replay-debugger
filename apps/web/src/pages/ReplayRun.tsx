import { useParams } from "react-router-dom";
import { useReplayStatus } from "@/hooks/useReplayStatus";
import { DivergenceReport } from "@/components/replay/DivergenceReport";

export function ProgressStepper({ status }: { status: string }) {
  const steps = ["pending", "running", "succeeded"];
  return (
    <ol className="progress-stepper">
      {steps.map((s) => (
        <li key={s} data-active={status === s}>
          {s}
        </li>
      ))}
    </ol>
  );
}

export function ReplayRunPage() {
  const { replayId = "" } = useParams();
  const { data } = useReplayStatus(replayId);
  const status = data?.status ?? "pending";

  return (
    <div>
      <h1>Replay {replayId.slice(0, 8)}</h1>
      <ProgressStepper status={status} />
      {status === "diverged" && (
        <DivergenceReport
          mismatchIndex={data?.divergence_index ?? 0}
          expected="success"
          actual="error"
        />
      )}
    </div>
  );
}

export default ReplayRunPage;
