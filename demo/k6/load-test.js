import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 1000,
  duration: '30s',
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<1000'],
  },
};

const targets = [
  'http://demo-api-production:8080/api',
  'http://demo-api-staging:8080/api',
];
const methods = ['GET', 'POST', 'PUT'];

export default function () {
  const target = targets[(__VU + __ITER) % targets.length];
  const method = methods[(__VU + __ITER) % methods.length];
  const response = http.request(method, target, method === 'POST' || method === 'PUT' ? JSON.stringify({ source: 'k6' }) : null, { headers: { 'Content-Type': 'application/json' } });
  check(response, { 'returns 200': (r) => r.status === 200 });
}
