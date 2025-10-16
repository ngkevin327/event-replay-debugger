import { useQuery } from "@tanstack/react-query";

export function useProjectTopics(projectId: string) {
  return useQuery({
    queryKey: ["topics", projectId],
    queryFn: async () => ["payments", "notifications", "ledger"],
    enabled: !!projectId,
  });
}
