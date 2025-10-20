import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/api/client";
import type { Incident } from "@/api/generated";

export function useIncidentStatus(incidentId: string) {
  return useQuery({
    queryKey: ["incident", incidentId],
    queryFn: () => apiFetch<Incident>(`/v1/incidents/${incidentId}`),
    refetchInterval: (q) => {
      const s = q.state.data?.status;
      if (s === "ready" || s === "failed") return false;
      return 2000;
    },
    enabled: !!incidentId,
  });
}
