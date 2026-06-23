import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';

const BASE_URL = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/$/, '');
const API = `${BASE_URL}/api`;

const USERNAME = __ENV.USERNAME || '';
const PASSWORD = __ENV.PASSWORD || '';
const PRODUCT_ID = Number(__ENV.PRODUCT_ID || 1);
const SKU_ID = Number(__ENV.SKU_ID || 1);
const ADDRESS_ID = Number(__ENV.ADDRESS_ID || 1);
const RUN_WRITES = (__ENV.RUN_WRITES || 'false').toLowerCase() === 'true';

const apiFailRate = new Rate('myshop_api_failed');
const businessFailCount = new Counter('myshop_business_failed');
const endpointDuration = new Trend('myshop_endpoint_duration', true);

export const options = {
  scenarios: {
    read_traffic: {
      executor: 'ramping-vus',
      stages: [
        { duration: __ENV.RAMP_UP || '30s', target: Number(__ENV.VUS || 50) },
        { duration: __ENV.HOLD || '1m', target: Number(__ENV.VUS || 50) },
        { duration: __ENV.RAMP_DOWN || '20s', target: 0 },
      ],
      gracefulRampDown: '10s',
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.05'],
    http_req_duration: ['p(95)<1000', 'p(99)<2000'],
    myshop_api_failed: ['rate<0.05'],
    'myshop_endpoint_duration{endpoint:product_detail}': ['p(95)<1000', 'p(99)<2000'],
    'myshop_endpoint_duration{endpoint:products}': ['p(95)<1000', 'p(99)<2000'],
    'myshop_endpoint_duration{endpoint:hot_products}': ['p(95)<1000', 'p(99)<2000'],
    'myshop_endpoint_duration{endpoint:seckill_list}': ['p(95)<1000', 'p(99)<2000'],
  },
};

function jsonHeaders(token) {
  const headers = {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  return { headers };
}

function parseJson(res) {
  try {
    return res.json();
  } catch (_) {
    return null;
  }
}

function checkApi(res, name, allowCodes = [200]) {
  const body = parseJson(res);
  const ok = check(res, {
    [`${name}: HTTP 200`]: (r) => r.status === 200,
    [`${name}: business code ok`]: () => body && allowCodes.includes(body.code),
  });

  apiFailRate.add(!ok);
  if (!ok) {
    businessFailCount.add(1);
  }

  return body;
}

function metricName(name) {
  return name.toLowerCase().replace(/[^a-z0-9]+/g, '_').replace(/^_|_$/g, '');
}

function recordEndpointDuration(res, name) {
  endpointDuration.add(res.timings.duration, { endpoint: metricName(name) });
}

function get(path, name, token) {
  const res = http.get(`${API}${path}`, jsonHeaders(token));
  recordEndpointDuration(res, name);
  return checkApi(res, name);
}

function post(path, payload, name, token, allowCodes = [200]) {
  const res = http.post(`${API}${path}`, JSON.stringify(payload), jsonHeaders(token));
  recordEndpointDuration(res, name);
  return checkApi(res, name, allowCodes);
}

export function setup() {
  if (!USERNAME || !PASSWORD) {
    return { token: '', authEnabled: false };
  }

  const body = post('/user/login', { username: USERNAME, password: PASSWORD }, 'login');
  const token = body && body.data && body.data.token ? body.data.token : '';

  if (!token) {
    throw new Error('Login failed. Check USERNAME/PASSWORD or run without auth endpoints.');
  }

  return { token, authEnabled: true };
}

export default function (data) {
  const token = data.token;

  group('public homepage and product APIs', () => {
    get('/index/banners', 'banners');
    get('/index/categories', 'categories');
    get('/index/products', 'products');
    get(`/index/product/detail?id=${PRODUCT_ID}`, 'product detail');
    get('/index/products/hot', 'hot products');
    get('/index/seckill/list?page=1&pageSize=10', 'seckill list');
  });

  if (data.authEnabled) {
    group('authenticated user APIs', () => {
      get('/auth/user/display/profile', 'user profile', token);
      get('/auth/cart/list', 'cart list', token);
      get('/auth/address/list', 'address list', token);
      get('/auth/order/list?status=-1&page=1&pageSize=10', 'order list', token);
    });

    if (RUN_WRITES) {
      group('optional write APIs', () => {
        post('/auth/cart/add', { sku_id: SKU_ID, quantity: 1 }, 'add cart', token);
        post('/auth/order/preview', { source_type: 'direct', sku_id: SKU_ID, quantity: 1 }, 'order preview', token, [200, 400]);
        post('/auth/order/create', {
          source_type: 'direct',
          sku_id: SKU_ID,
          quantity: 1,
          address_id: ADDRESS_ID,
          remark: 'k6 load test',
        }, 'create order', token, [200, 400]);
      });
    }
  }

  sleep(Number(__ENV.SLEEP || 1));
}
