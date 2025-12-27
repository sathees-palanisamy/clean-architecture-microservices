import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Counter } from 'k6/metrics';

const stockOutCount = new Counter('stock_out_count');

export const options = {
  stages: [
    { duration: '30s', target: 20 },
    { duration: '1m', target: 20 },
    { duration: '20s', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
  },
};

const BASE_URL_PRODUCT = __ENV.PRODUCT_SERVICE_URL || 'http://localhost:8081';
const BASE_URL_ORDER = __ENV.ORDER_SERVICE_URL || 'http://localhost:8082';

export default function () {
  group('Product Service', () => {
    // 1. Health Check
    let healthRes = http.get(`${BASE_URL_PRODUCT}/health`);
    check(healthRes, { 'product health status is 200': (r) => r.status === 200 });

    // 2. Get All Products
    let productsRes = http.get(`${BASE_URL_PRODUCT}/products`);
    check(productsRes, { 'get products status is 200': (r) => r.status === 200 });

    let products = productsRes.json();
    if (products && products.length > 0) {
      let product = products[Math.floor(Math.random() * products.length)];

      // 3. Get Single Product
      let singleProductRes = http.get(`${BASE_URL_PRODUCT}/products/${product.id}`);
      check(singleProductRes, { 'get single product status is 200': (r) => r.status === 200 });

      // 4. Stock Reserve (Internal logic, but exposed via API)
      let reservePayload = JSON.stringify({ product_id: product.id, quantity: 1 });
      let reserveRes = http.post(`${BASE_URL_PRODUCT}/products/reserve`, reservePayload, { headers: { 'Content-Type': 'application/json' } });
      check(reserveRes, { 'reserve stock status is 200 or 422': (r) => r.status === 200 || r.status === 422 });
    }

    // 5. Create Product (Simulate admin behavior occasionally)
    if (Math.random() < 0.1) {
      let createPayload = JSON.stringify({
        sku: `SKU-${Date.now()}-${Math.floor(Math.random() * 1000)}`,
        name: 'Load Test Product',
        price: 10.50,
        total_qty: 100
      });
      let createRes = http.post(`${BASE_URL_PRODUCT}/products`, createPayload, { headers: { 'Content-Type': 'application/json' } });
      check(createRes, { 'create product status is 201': (r) => r.status === 201 });
    }
  });

  group('Order Service', () => {
    // 1. Health Check
    let healthRes = http.get(`${BASE_URL_ORDER}/health`);
    check(healthRes, { 'order health status is 200': (r) => r.status === 200 });

    // 2. Create Order
    let productsRes = http.get(`${BASE_URL_PRODUCT}/products`);
    let products = productsRes.json();
    if (products && products.length > 0) {
      let product = products[Math.floor(Math.random() * products.length)];
      let orderPayload = JSON.stringify({
        user_id: Math.floor(Math.random() * 1000) + 1,
        product_id: product.id,
        quantity: 1,
      });

      let orderRes = http.post(`${BASE_URL_ORDER}/orders`, orderPayload, { headers: { 'Content-Type': 'application/json' } });
      if (orderRes.status === 422) stockOutCount.add(1);
      check(orderRes, {
        'order status is successful (20x or 422)': (r) => r.status === 201 || r.status === 422,
        'no server error': (r) => r.status < 500,
      });
    }

    // 3. Get All Orders
    let allOrdersRes = http.get(`${BASE_URL_ORDER}/orders`);
    check(allOrdersRes, { 'get all orders status is 200': (r) => r.status === 200 });

    let orders = allOrdersRes.json();
    if (orders && orders.length > 0) {
      let order = orders[Math.floor(Math.random() * orders.length)];
      // 4. Get Single Order
      let singleOrderRes = http.get(`${BASE_URL_ORDER}/orders/${order.id}`);
      check(singleOrderRes, { 'get single order status is 200': (r) => r.status === 200 });
    }
  });

  sleep(1);
}
