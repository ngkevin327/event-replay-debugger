import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "@/api/client";
import type { ReplayRun } from "@/api/generated";

export function useReplayStatus(replayId: string) {
  return useQuery({
    queryKey: ["replay", replayId],
    queryFn: () =>
      apiFetch<{ replay: ReplayRun }>(`/v1/replays/${replayId}`).then(
        (r) => r.replay ?? (r as unknown as ReplayRun),
      ),
    refetchInterval: (q) => {
      const s = q.state.data?.status;
      if (s === "succeeded" || s === "failed" || s === "diverged" || s === "cancelled")
        return false;
      return 1500;
    },
    enabled: !!replayId,
  });
}
