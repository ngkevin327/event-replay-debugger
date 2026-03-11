import { JsonViewer } from "@/components/JsonViewer";
import type { TimelineEvent } from "./types";

export function RetryBadge({ generation }: { generation: number }) {
  return <span className="retry-badge">Retry gen {generation}</span>;
}

export function EventDrawer({
  event,
  open,
  onClose,
}: {
  event: TimelineEvent | null;
  open: boolean;
  onClose: () => void;
}) {
  if (!open || !event) return null;
  return (
    <>
      <div
        className="event-drawer-backdrop"
        role="presentation"
        onClick={onClose}
        onKeyDown={() => {}}
      />
      <aside className="event-drawer" role="dialog" aria-label="Event details">
        <div className="event-drawer__header">
          <div>
            <h3>{event.event_id}</h3>
            {event.retry_generation != null && (
              <RetryBadge generation={event.retry_generation} />
            )}
          </div>
          <button type="button" className="btn btn--ghost" onClick={onClose}>
            Close
          </button>
        </div>
        <p className="timeline-row__meta mb-4">
          {event.topic} · partition {event.partition} · offset {event.offset}
        </p>
        <JsonViewer data={event} />
      </aside>
    </>
  );
}
