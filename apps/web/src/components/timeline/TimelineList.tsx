import { memo, useCallback } from "react";
import { FixedSizeList as List, type ListChildComponentProps } from "react-window";
import type { TimelineEvent } from "./types";
import { GapMarker } from "./GapMarker";
import { RetryGroup } from "./RetryGroup";

export type TimelineRowItem =
  | { type: "event"; event: TimelineEvent & { rowKey?: string } }
  | { type: "gap"; gap: { topic: string; partition: number; start_offset: number; end_offset: number } }
  | { type: "group"; label: string };

export function groupRows(events: TimelineEvent[]): TimelineRowItem[] {
  return events.map((e) => ({ type: "event" as const, event: e }));
}

type RowProps = ListChildComponentProps<TimelineRowItem[]>;

const Row = memo(function Row({ index, style, data }: RowProps) {
  const row = data[index];
  if (row.type === "gap") {
    return (
      <div style={style}>
        <GapMarker gap={row.gap} />
      </div>
    );
  }
  if (row.type === "group") {
    return (
      <div style={style}>
        <RetryGroup label={row.label} />
      </div>
    );
  }
  const e = row.event;
  return (
    <div style={style} className="timeline-row" data-event-id={e.event_id}>
      <span>{e.topic}</span> p{e.partition} @{e.offset}
      {e.retry_generation ? ` gen${e.retry_generation}` : ""}
    </div>
  );
});

export function TimelineList({
  rows,
  height = 400,
  activeIndex,
  onActiveChange,
}: {
  rows: TimelineRowItem[];
  height?: number;
  activeIndex?: number;
  onActiveChange?: (i: number) => void;
}) {
  const onItemsRendered = useCallback(() => {
    if (activeIndex != null) onActiveChange?.(activeIndex);
  }, [activeIndex, onActiveChange]);

  return (
    <List
      height={height}
      width="100%"
      itemCount={rows.length}
      itemSize={36}
      itemData={rows}
      onItemsRendered={onItemsRendered}
    >
      {Row}
    </List>
  );
}
