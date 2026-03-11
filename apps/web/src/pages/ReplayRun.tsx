import { useParams, Link } from "react-router-dom";
import { useReplayStatus } from "@/hooks/useReplayStatus";
import { DivergenceReport } from "@/components/replay/DivergenceReport";
import { PageHeader } from "@/components/PageHeader";
import { StatusBadge } from "@/components/StatusBadge";

const STEPS = ["pending", "running", "succeeded"] as const;

export function ProgressStepper({ status }: { status: string }) {
  const idx = STEPS.indexOf(status as (typeof STEPS)[number]);
  return (
    <ol className="progress-stepper" aria-label="Replay progress">
      {STEPS.map((s, i) => (
        <li key={s} data-active={status === s} data-done={idx > i}>
          <StatusBadge status={s} />
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
      <PageHeader
        title={`Replay ${replayId.slice(0, 8)}…`}
        description="Track replay execution and review divergence when outcomes differ."
        actions={
          <>
            <StatusBadge status={status} />
            <Link to="/incidents" className="btn btn--secondary">
              Back to incidents
            </Link>
          </>
        }
      />
      <ProgressStepper status={status} />
      {status === "diverged" && (
        <div className="divergence-card animate-in">
          <DivergenceReport
            mismatchIndex={data?.divergence_index ?? 0}
            expected="success"
            actual="error"
          />
        </div>
      )}
      {status === "pending" && (
        <p className="empty-state mt-6">
          Replay is queued. The orchestrator will start the worker shortly.
        </p>
      )}
    </div>
  );
}

export default ReplayRunPage;
