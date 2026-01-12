import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  scenarios: {
    sustained_ingest: {
      executor: "constant-arrival-rate",
      rate: 10000,
      timeUnit: "1s",
      duration: "2m",
      preAllocatedVUs: 50,
      maxVUs: 200,
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<500"],
  },
};

const target = __ENV.K6_TARGET_URL || "http://localhost:8081";
const apiKey = __ENV.K6_API_KEY || "test-key";

export default function () {
  const payload = JSON.stringify({
    events: [
      {
        event_id: `e-${__VU}-${__ITER}`,
        project_id: "00000000-0000-0000-0000-000000000099",
        captured_at: new Date().toISOString(),
        source: "consumer",
        topic: "payments",
        partition: 0,
        offset: __ITER,
        timestamp: new Date().toISOString(),
        outcome: "success",
      },
    ],
  });
  const res = http.post(`${target}/v1/ingest/batch`, payload, {
    headers: {
      "Content-Type": "application/json",
      "X-Replay-Key": apiKey,
    },
  });
  check(res, { "accepted": (r) => r.status === 202 || r.status === 200 });
  sleep(0.01);
}
