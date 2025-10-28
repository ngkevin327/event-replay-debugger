import { JsonViewer } from "@/components/JsonViewer";
import type { TimelineEvent } from "./types";

export function RetryBadge({ generation }: { generation: number }) {
  return <span className="retry-badge">retry gen {generation}</span>;
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
    <aside className="event-drawer" role="dialog" aria-label="Event details">
      <button type="button" onClick={onClose}>
        Close
      </button>
      <h3>{event.event_id}</h3>
      {event.retry_generation != null && (
        <RetryBadge generation={event.retry_generation} />
      )}
      <p>
        {event.topic} p{event.partition} @{event.offset}
      </p>
      <JsonViewer data={event} />
    </aside>
  );
}
