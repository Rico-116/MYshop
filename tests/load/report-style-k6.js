import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';

const BASE_URL = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/$/, '');
const API = `${BASE_URL}/api`;

const MODE = __ENV.MODE || 'hotspot';
const VUS = Number(__ENV.VUS || 100);
const DURATION = __ENV.DURATION || '60s';
const PRODUCT_ID = Number(__ENV.PRODUCT_ID || 1);
const CATEGORY_ID = Number(__ENV.CATEGORY_ID || 1);
const KEYWORD = __ENV.KEYWORD || '';

const apiFailRate = new Rate('myshop_api_failed');
const businessFailCount = new Counter('myshop_business_failed');
const endpointDuration = new Trend('myshop_endpoint_duration', true);

export const options = {
  scenarios: {
    report_style: {
      executor: 'constant-vus',
      vus: VUS,
      duration: DURATION,
      gracefulStop: '30s',
    },
  },
  summaryTrendStats: ['avg', 'min', 'med', 'p(90)', 'p(95)', 'p(99)', 'max'],
  thresholds: {
    http_req_failed: ['rate<0.05'],
    myshop_api_failed: ['rate<0.05'],
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
    tags: { endpoint: metricName(name), mode: MODE },
  });

  endpointDuration.add(res.timings.duration, { endpoint: metricName(name), mode: MODE });

  const body = parseJson(res);
  const ok = check(res, {
    [`${name}: HTTP 200`]: (r) => r.status === 200,
    [`${name}: business code 200`]: () => body && body.code === 200,
  });

  apiFailRate.add(!ok);
  if (!ok) {
    businessFailCount.add(1);
  }
}

function hotspot() {
  request(`/index/product/detail?id=${PRODUCT_ID}`, 'product detail');
}

function mixedPublicRead() {
  const n = Math.random();

  if (n < 0.70) {
    request(`/index/product/detail?id=${PRODUCT_ID}`, 'product detail');
    return;
  }

  if (n < 0.85) {
    request('/index/products', 'products');
    return;
  }

  if (n < 0.95) {
    request(`/index/search?keyword=${encodeURIComponent(KEYWORD)}&page=1&pageSize=10`, 'search');
    return;
  }

  request(`/index/category/display?category_id=${CATEGORY_ID}`, 'category display');
}

export default function () {
  if (MODE === 'mix') {
    mixedPublicRead();
  } else {
    hotspot();
  }

  sleep(0);
}

export function handleSummary(data) {
  const output = {};

  if (__ENV.SUMMARY_FILE) {
    output[__ENV.SUMMARY_FILE] = JSON.stringify(data, null, 2);
  } else {
    output.stdout = JSON.stringify(data, null, 2);
  }

  return output;
}
