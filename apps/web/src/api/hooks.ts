import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiFetch } from "./client";
import type { Agent, Incident, Project, ReplayRun } from "./generated";

export function useProjects() {
  return useQuery({
    queryKey: ["projects"],
    queryFn: () => apiFetch<{ projects: Project[] }>("/v1/projects"),
  });
}

export function useIncidents(projectId: string, status?: string) {
  return useQuery({
    queryKey: ["incidents", projectId, status],
    queryFn: () => {
      const q = status ? `?status=${status}` : "";
      return apiFetch<{ incidents: Incident[] }>(
        `/v1/projects/${projectId}/incidents${q}`,
      );
    },
    enabled: !!projectId,
  });
}

export function useAgentHealth() {
  return useQuery({
    queryKey: ["agents"],
    queryFn: () => apiFetch<{ agents: Agent[] }>("/v1/agents"),
  });
}

export function useRecentIncidents(projectId: string) {
  return useIncidents(projectId);
}

export function useApiKeys(_projectId: string) {
  return useQuery({
    queryKey: ["api-keys", _projectId],
    queryFn: async () => ({ keys: [] as { id: string; prefix: string }[] }),
  });
}

export function useUpdateAllowlist() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (_topics: string[]) => undefined,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["projects"] }),
  });
}

export function useProjectUsage(_projectId: string) {
  return useQuery({
    queryKey: ["usage", _projectId],
    queryFn: async () => ({ events_used: 0, events_limit: 1_000_000, percent: 0 }),
  });
}

export function useCreateReplay(incidentId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (timing_mode: string) =>
      apiFetch<ReplayRun>(`/v1/incidents/${incidentId}/replays`, {
        method: "POST",
        body: JSON.stringify({ timing_mode }),
      }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["replays", incidentId] }),
  });
}

export function useReplayList(_incidentId: string) {
  return useQuery({
    queryKey: ["replays", _incidentId],
    queryFn: async () => ({ replays: [] as ReplayRun[] }),
  });
}

export function useExportIncident(incidentId: string) {
  return useMutation({
    mutationFn: async () => {
      const res = await fetch(`/v1/incidents/${incidentId}/export`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("replay_token") ?? ""}`,
        },
      });
      if (!res.ok) throw new Error("export failed");
      const blob = await res.blob();
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `incident-${incidentId}.json`;
      a.click();
      URL.revokeObjectURL(url);
    },
  });
}

export function useCreateShare(incidentId: string) {
  return useMutation({
    mutationFn: (ttl_hours: number) =>
      apiFetch<{ token: string; expires_at: string; url: string }>(
        `/v1/incidents/${incidentId}/share-tokens`,
        { method: "POST", body: JSON.stringify({ ttl_hours }) },
      ),
  });
}
