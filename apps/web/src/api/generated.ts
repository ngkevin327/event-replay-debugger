/** OpenAPI-derived path and schema types (MVP subset). */

export type IncidentStatus = "collecting" | "ready" | "failed";
export type ReplayStatus =
  | "pending"
  | "running"
  | "succeeded"
  | "failed"
  | "diverged"
  | "cancelled";

export interface Project {
  id: string;
  name: string;
}

export interface Incident {
  id: string;
  project_id?: string;
  status: IncidentStatus;
  window_start?: string;
  window_end?: string;
  topic_filters?: string[];
  event_count?: number;
  coverage_percent?: number;
}

export interface Agent {
  id: string;
  hostname?: string;
  status?: string;
  last_heartbeat?: string;
}

export interface ReplayRun {
  id: string;
  incident_id?: string;
  status: ReplayStatus;
  timing_mode?: string;
  divergence_index?: number | null;
}

export interface TimelineResponse {
  version: number;
  timeline: { events?: TimelineEvent[]; gaps?: GapRange[] };
}

export interface TimelineEvent {
  event_id: string;
  topic: string;
  partition: number;
  offset: number;
  arrival_index?: number;
  retry_generation?: number;
  outcome?: string;
  correlation_id?: string;
  payload_hash?: string;
}

export interface GapRange {
  topic: string;
  partition: number;
  start_offset: number;
  end_offset: number;
}

export interface GraphPayload {
  nodes: { id: string; kind: string; topic?: string; failed?: boolean }[];
  edges: { from: string; to: string }[];
}

export type Paths = {
  "/v1/projects": { get: { responses: { 200: { projects: Project[] } } } };
  "/v1/projects/{projectId}/incidents": {
    get: { responses: { 200: { incidents: Incident[] } } };
    post: { body: { window_start: string; window_end: string; topic_filters?: string[] } };
  };
  "/v1/incidents/{incidentId}": { get: { responses: { 200: Incident } } };
  "/v1/incidents/{incidentId}/timeline": { get: { responses: { 200: TimelineResponse } } };
  "/v1/incidents/{incidentId}/graph": { get: { responses: { 200: GraphPayload } } };
  "/v1/incidents/{incidentId}/replays": {
    post: { body: { timing_mode?: string }; responses: { 201: ReplayRun } };
  };
  "/v1/replays/{replayId}": { get: { responses: { 200: ReplayRun } }; delete: unknown };
  "/v1/agents": { get: { responses: { 200: { agents: Agent[] } } } };
};
