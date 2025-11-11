import { useState, useCallback } from "react";

export function useIncidentSelection() {
  const [selectedEventId, setSelectedEventId] = useState<string | null>(null);
  const select = useCallback((id: string | null) => setSelectedEventId(id), []);
  return { selectedEventId, select };
}
