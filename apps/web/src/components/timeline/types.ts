export type { TimelineEvent, GapRange } from "@/api/generated";

export type TimelineRow =
  | { kind: "event"; id: string; data: import("@/api/generated").TimelineEvent }
  | { kind: "gap"; id: string; data: import("@/api/generated").GapRange }
  | { kind: "retry-group"; id: string; label: string };
