import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';

const BASE_URL = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/$/, '');
const API = `${BASE_URL}/api`;

const PRODUCT_ID = Number(__ENV.PRODUCT_ID || 1);
const TARGET_QPS = Number(__ENV.TARGET_QPS || 100);
const PRE_ALLOCATED_VUS = Number(__ENV.PRE_ALLOCATED_VUS || Math.max(200, TARGET_QPS * 5));
const MAX_VUS = Number(__ENV.MAX_VUS || Math.max(1000, TARGET_QPS * 20));
const DURATION = __ENV.DURATION || '2m';
const ENDPOINT = __ENV.ENDPOINT || 'mix';

const apiFailRate = new Rate('myshop_api_failed');
const businessFailCount = new Counter('myshop_business_failed');
const endpointDuration = new Trend('myshop_endpoint_duration', true);

export const options = {
  scenarios: {
    max_qps_probe: {
      executor: 'constant-arrival-rate',
      rate: TARGET_QPS,
      timeUnit: '1s',
      duration: DURATION,
      preAllocatedVUs: PRE_ALLOCATED_VUS,
      maxVUs: MAX_VUS,
      gracefulStop: '30s',
    },
  },
  summaryTrendStats: ['avg', 'min', 'med', 'p(90)', 'p(95)', 'p(99)', 'p(99.9)', 'max'],
  thresholds: {
    http_req_failed: ['rate<0.01'],
    myshop_api_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<1000', 'p(99)<2000'],
  },
};

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

function request(path, name) {
  const res = http.get(`${API}${path}`, {
    headers: { Accept: 'application/json' },
  });

  endpointDuration.add(res.timings.duration, { endpoint: metricName(name) });

  const body = parseJson(res);
  const ok = check(res, {
    [`${name}: HTTP 200`]: (r) => r.status === 200,
    [`${name}: business code ok`]: () => body && body.code === 200,
  });

  apiFailRate.add(!ok);
  if (!ok) businessFailCount.add(1);
}

function runOneEndpoint() {
  switch (ENDPOINT) {
    case 'products':
      request('/index/products', 'products');
      break;
    case 'product_detail':
      request(`/index/product/detail?id=${PRODUCT_ID}`, 'product detail');
      break;
    case 'hot_products':
      request('/index/products/hot', 'hot products');
      break;
    case 'seckill_list':
      request('/index/seckill/list?page=1&pageSize=10', 'seckill list');
      break;
    case 'categories':
      request('/index/categories', 'categories');
      break;
    case 'banners':
      request('/index/banners', 'banners');
      break;
    default:
      {
        const pick = Math.floor(Math.random() * 4);
        if (pick === 0) request('/index/products', 'products');
        if (pick === 1) request(`/index/product/detail?id=${PRODUCT_ID}`, 'product detail');
        if (pick === 2) request('/index/products/hot', 'hot products');
        if (pick === 3) request('/index/seckill/list?page=1&pageSize=10', 'seckill list');
      }
  }
}

export default function () {
  runOneEndpoint();
  sleep(0);
}

export function handleSummary(data) {
  const output = {
    stdout: JSON.stringify({
      target_qps: TARGET_QPS,
      endpoint: ENDPOINT,
      duration: DURATION,
      base_url: BASE_URL,
      metrics: data.metrics,
      root_group: data.root_group,
    }, null, 2),
  };

  if (__ENV.SUMMARY_FILE) {
    output[__ENV.SUMMARY_FILE] = JSON.stringify(data, null, 2);
  }

  return output;
}
