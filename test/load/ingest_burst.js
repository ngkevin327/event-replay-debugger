import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    burst: {
      executor: 'constant-arrival-rate',
      rate: 10000,
      timeUnit: '1s',
      duration: '10s',
      preAllocatedVUs: 50,
      maxVUs: 200,
    },
  },
};

const url = __ENV.INGEST_URL || 'http://localhost:8081/v1/ingest/batch';
const apiKey = __ENV.REPLAY_API_KEY || '';

export default function () {
  const payload = JSON.stringify({
    events: [
      {
        event_id: '11111111-1111-1111-1111-111111111111',
        project_id: __ENV.PROJECT_ID || '00000000-0000-0000-0000-000000000002',
        captured_at: new Date().toISOString(),
        source: 'consumer',
        topic: 'orders',
        partition: 0,
        offset: Math.floor(Math.random() * 1000000),
        timestamp: new Date().toISOString(),
        payload: '{"demo":true}',
      },
    ],
  });
  const res = http.post(url, payload, {
    headers: {
      'Content-Type': 'application/json',
      'X-Replay-Key': apiKey,
    },
  });
  check(res, { 'status is 202 or 429': (r) => r.status === 202 || r.status === 429 });
  sleep(0.01);
}
