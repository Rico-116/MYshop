import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';

const BASE_URL = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/$/, '');
const API = `${BASE_URL}/api`;

const USERNAME = __ENV.USERNAME || '';
const PASSWORD = __ENV.PASSWORD || '';
const PRODUCT_ID = Number(__ENV.PRODUCT_ID || 1);
const TARGET_VUS = Number(__ENV.VUS || 50);

const apiFailRate = new Rate('myshop_api_failed');
const businessFailCount = new Counter('myshop_business_failed');
const endpointDuration = new Trend('myshop_endpoint_duration', true);

export const options = {
  scenarios: {
    tail_latency: {
      executor: 'constant-vus',
      vus: TARGET_VUS,
      duration: __ENV.DURATION || '10m',
      gracefulStop: '30s',
    },
  },
  summaryTrendStats: ['avg', 'min', 'med', 'p(90)', 'p(95)', 'p(99)', 'p(99.9)', 'max'],
  thresholds: {
    http_req_failed: ['rate<0.01'],
    myshop_api_failed: ['rate<0.01'],
    http_req_duration: [
      `p(95)<${Number(__ENV.P95_MS || 500)}`,
      `p(99)<${Number(__ENV.P99_MS || 1000)}`,
      `p(99.9)<${Number(__ENV.P999_MS || 3000)}`,
    ],
  },
};

function headers(token) {
  const h = { Accept: 'application/json', 'Content-Type': 'application/json' };
  if (token) h.Authorization = `Bearer ${token}`;
  return { headers: h };
}

function parseJson(res) {
  try {
    return res.json();
  } catch (_) {
    return null;
  }
}

function metricName(name) {
  return name.toLowerCase().replace(/[^a-z0-9]+/g, '_').replace(/^_|_$/g, '');
}

function checkApi(res, name, allowCodes = [200]) {
  endpointDuration.add(res.timings.duration, { endpoint: metricName(name) });
  const body = parseJson(res);
  const ok = check(res, {
    [`${name}: HTTP 200`]: (r) => r.status === 200,
    [`${name}: business code ok`]: () => body && allowCodes.includes(body.code),
  });

  apiFailRate.add(!ok);
  if (!ok) businessFailCount.add(1);

  return body;
}

function get(path, name, token) {
  const res = http.get(`${API}${path}`, headers(token));
  return checkApi(res, name);
}

function post(path, payload, name, token) {
  const res = http.post(`${API}${path}`, JSON.stringify(payload), headers(token));
  return checkApi(res, name);
}

export function setup() {
  if (!USERNAME || !PASSWORD) return { token: '', authEnabled: false };

  const body = post('/user/login', { username: USERNAME, password: PASSWORD }, 'login');
  const token = body && body.data && body.data.token ? body.data.token : '';
  if (!token) throw new Error('Login failed. Check USERNAME/PASSWORD.');

  return { token, authEnabled: true };
}

export default function (data) {
  group('tail latency read path', () => {
    get('/index/products', 'products');
    get(`/index/product/detail?id=${PRODUCT_ID}`, 'product detail');
    get('/index/products/hot', 'hot products');
    get('/index/seckill/list?page=1&pageSize=10', 'seckill list');
  });

  if (data.authEnabled) {
    group('tail latency auth read path', () => {
      get('/auth/cart/list', 'cart list', data.token);
      get('/auth/order/list?status=-1&page=1&pageSize=10', 'order list', data.token);
    });
  }

  sleep(Number(__ENV.SLEEP || 0.2));
}
