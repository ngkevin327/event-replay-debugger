import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/api/client";
import type { TimelineEvent, GapRange } from "@/api/generated";

export function useTimelineEvents(incidentId: string) {
  const q = useQuery({
    queryKey: ["timeline", incidentId],
    queryFn: () =>
      apiFetch<{ version: number; timeline: { events?: TimelineEvent[]; gaps?: GapRange[] } }>(
        `/v1/incidents/${incidentId}/timeline`,
      ),
    enabled: !!incidentId,
  });

  const events = useMemo(() => {
    const raw = q.data?.timeline?.events ?? [];
    return raw.map((e, i) => ({
      ...e,
      rowKey: `${e.topic}-${e.partition}-${e.offset}-${i}`,
    }));
  }, [q.data]);

  return { ...q, events, gaps: q.data?.timeline?.gaps ?? [] };
}
