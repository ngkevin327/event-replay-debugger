import { useState, useMemo } from "react";
import { useParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/api/client";
import { CoverageBar } from "@/components/CoverageBar";
import { StatusBadge } from "@/components/StatusBadge";
import {
  ReconstructionProgress,
  SkeletonTimeline,
} from "@/components/LoadingStates";
import { useIncidentStatus } from "@/hooks/useIncidentStatus";
import { useTimelineEvents } from "@/hooks/useTimelineEvents";
import { TimelineList, groupRows, type TimelineRowItem } from "@/components/timeline/TimelineList";
import { TimelineScrubber } from "@/components/timeline/TimelineScrubber";
import { useTimelineKeyboard } from "@/hooks/useTimelineKeyboard";
import { EventDrawer } from "@/components/timeline/EventDrawer";
import { SnapshotsPanel } from "@/components/timeline/SnapshotsPanel";
import { WorkflowGraph } from "@/components/graph/WorkflowGraph";
import { useIncidentSelection } from "@/hooks/useIncidentSelection";
import { ReplayPanel } from "@/components/replay/ReplayPanel";
import { ReplayHistory } from "@/components/replay/ReplayHistory";
import { IncidentActions } from "@/components/IncidentActions";
import { useReplayList } from "@/api/hooks";
import type { GraphPayload } from "@/api/generated";
import type { TimelineEvent } from "@/components/timeline/types";

export function IncidentHeader({
  status,
  coverage,
}: {
  status: string;
  coverage: number;
}) {
  return (
    <header className="incident-header">
      <h1>Incident</h1>
      <div className="incident-header__meta">
        <StatusBadge status={status} />
        <CoverageBar percent={coverage} />
      </div>
    </header>
  );
}

export default function IncidentDetail() {
  const { incidentId = "" } = useParams();
  const { data: incident, isLoading } = useIncidentStatus(incidentId);
  const timeline = useTimelineEvents(incidentId);
  const graphQ = useQuery({
    queryKey: ["graph", incidentId],
    queryFn: () => apiFetch<GraphPayload>(`/v1/incidents/${incidentId}/graph`),
    enabled: incident?.status === "ready",
  });
  const replays = useReplayList(incidentId);
  const { select } = useIncidentSelection();
  const [tab, setTab] = useState<"timeline" | "snapshots" | "replay">("timeline");
  const [activeIndex, setActiveIndex] = useState(0);
  const [drawerOpen, setDrawerOpen] = useState(false);

  const rows: TimelineRowItem[] = useMemo(() => {
    const events = timeline.events ?? [];
    const gapRows: TimelineRowItem[] = (timeline.gaps ?? []).map((g) => ({
      type: "gap" as const,
      gap: g,
    }));
    return [...groupRows(events), ...gapRows];
  }, [timeline.events, timeline.gaps]);

  useTimelineKeyboard(activeIndex, rows.length, setActiveIndex);

  const selectedEvent: TimelineEvent | null =
    rows[activeIndex]?.type === "event" ? rows[activeIndex].event : null;

  if (isLoading) return <p className="loading-state">Loading incident…</p>;
  if (!incident) return <p className="loading-state">Incident not found</p>;

  const ready = incident.status === "ready";

  return (
    <div className="incident-detail">
      <IncidentHeader
        status={incident.status}
        coverage={incident.coverage_percent ?? 0}
      />
      <IncidentActions incidentId={incidentId} />
      {!ready && (
        <>
          <ReconstructionProgress step={incident.status} />
          <SkeletonTimeline />
        </>
      )}
      {ready && (
        <>
          <div className="tabs" role="tablist">
            <button
              type="button"
              role="tab"
              className={tab === "timeline" ? "active" : ""}
              aria-selected={tab === "timeline"}
              onClick={() => setTab("timeline")}
            >
              Timeline
            </button>
            <button
              type="button"
              role="tab"
              className={tab === "snapshots" ? "active" : ""}
              aria-selected={tab === "snapshots"}
              onClick={() => setTab("snapshots")}
            >
              Snapshots
            </button>
            <button
              type="button"
              role="tab"
              className={tab === "replay" ? "active" : ""}
              aria-selected={tab === "replay"}
              onClick={() => setTab("replay")}
            >
              Replay
            </button>
          </div>
          {tab === "timeline" && (
            <section className="timeline-panel">
              <TimelineScrubber
                max={Math.max(0, rows.length - 1)}
                value={activeIndex}
                onChange={setActiveIndex}
              />
              <TimelineList
                rows={rows}
                activeIndex={activeIndex}
                onActiveChange={setActiveIndex}
              />
              {graphQ.data && (
                <WorkflowGraph
                  graph={graphQ.data}
                  onSelectNode={(id) => {
                    select(id);
                    document
                      .querySelector(`[data-event-id="${id}"]`)
                      ?.scrollIntoView({ block: "center" });
                  }}
                />
              )}
              <button
                type="button"
                className="btn btn--secondary mt-4"
                onClick={() => setDrawerOpen(true)}
                onKeyDown={(e) => e.key === "Enter" && setDrawerOpen(true)}
              >
                Event details
              </button>
              <EventDrawer
                event={selectedEvent}
                open={drawerOpen}
                onClose={() => setDrawerOpen(false)}
              />
            </section>
          )}
          {tab === "snapshots" && (
            <SnapshotsPanel
              snapshots={[
                {
                  consumer_group: "payment-cg",
                  topic: "payments",
                  partition: 0,
                  offset: 120,
                },
              ]}
            />
          )}
          {tab === "replay" && (
            <section className="panel-stack">
              <ReplayPanel incidentId={incidentId} />
              <ReplayHistory replays={replays.data?.replays ?? []} />
            </section>
          )}
        </>
      )}
    </div>
  );
}
