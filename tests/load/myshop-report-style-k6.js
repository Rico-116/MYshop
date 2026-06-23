import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';

const BASE_URL = (__ENV.BASE_URL || 'http://localhost:8080').replace(/\/$/, '');
const API = `${BASE_URL}/api`;

const MODE = (__ENV.MODE || 'hotspot').toLowerCase();
const VUS = Number(__ENV.VUS || 100);
const DURATION = __ENV.DURATION || '60s';
const PRODUCT_ID = Number(__ENV.PRODUCT_ID || 1);
const CATEGORY_ID = Number(__ENV.CATEGORY_ID || 1);
const KEYWORD = __ENV.KEYWORD || '';
const PAGE_SIZE = Number(__ENV.PAGE_SIZE || 10);
const THINK_TIME = Number(__ENV.THINK_TIME || 0);

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

function get(path, name) {
  const endpoint = metricName(name);
  const res = http.get(`${API}${path}`, {
    headers: { Accept: 'application/json' },
    tags: { scene: MODE, endpoint },
  });

  endpointDuration.add(res.timings.duration, { scene: MODE, endpoint });

  const body = parseJson(res);
  const ok = check(res, {
    [`${name}: http 200`]: (r) => r.status === 200,
    [`${name}: business code 200`]: () => body && body.code === 200,
  });

  apiFailRate.add(!ok, { scene: MODE, endpoint });
  if (!ok) {
    businessFailCount.add(1, { scene: MODE, endpoint });
  }

  return res;
}

function hotspotProductDetail() {
  get(`/index/product/detail?id=${PRODUCT_ID}`, 'product detail');
}

function cachedProductList() {
  get('/index/products', 'product list');
}

function hotProductList() {
  get('/index/products/hot', 'hot products');
}

function categoryDisplay() {
  get(`/index/category/display?category_id=${CATEGORY_ID}`, 'category display');
}

function publicReadMix() {
  const n = Math.random();

  if (n < 0.70) {
    get(`/index/product/detail?id=${PRODUCT_ID}`, 'product detail');
    return;
  }

  if (n < 0.85) {
    get('/index/products', 'product list');
    return;
  }

  if (n < 0.95) {
    get(`/index/search?keyword=${encodeURIComponent(KEYWORD)}&page=1&pageSize=${PAGE_SIZE}`, 'product search');
    return;
  }

  get(`/index/category/display?category_id=${CATEGORY_ID}`, 'category display');
}

export default function () {
  switch (MODE) {
    case 'mix':
      publicReadMix();
      break;
    case 'products':
      cachedProductList();
      break;
    case 'hot':
      hotProductList();
      break;
    case 'category':
      categoryDisplay();
      break;
    default:
      hotspotProductDetail();
      break;
  }

  if (THINK_TIME > 0) {
    sleep(THINK_TIME);
  }
}

export function handleSummary(data) {
  if (__ENV.SUMMARY_FILE) {
    return {
      [__ENV.SUMMARY_FILE]: JSON.stringify(data, null, 2),
      stdout: JSON.stringify(data.metrics, null, 2),
    };
  }

  return {
    stdout: JSON.stringify(data, null, 2),
  };
}
