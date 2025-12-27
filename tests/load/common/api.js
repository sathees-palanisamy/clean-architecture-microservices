import http from 'k6/http';
import { check } from 'k6';

const BASE_URL_PRODUCT = __ENV.PRODUCT_SERVICE_URL || 'http://localhost:8081';
const BASE_URL_ORDER = __ENV.ORDER_SERVICE_URL || 'http://localhost:8082';

export const Api = {
    product: {
        health: () => http.get(`${BASE_URL_PRODUCT}/health`),
        getAll: () => http.get(`${BASE_URL_PRODUCT}/products`),
        getOne: (id) => http.get(`${BASE_URL_PRODUCT}/products/${id}`),
        create: (data) => http.post(`${BASE_URL_PRODUCT}/products`, JSON.stringify(data), { headers: { 'Content-Type': 'application/json' } }),
        reserve: (id, qty) => http.post(`${BASE_URL_PRODUCT}/products/reserve`, JSON.stringify({ product_id: id, quantity: qty }), { headers: { 'Content-Type': 'application/json' } }),
    },
    order: {
        health: () => http.get(`${BASE_URL_ORDER}/health`),
        getAll: () => http.get(`${BASE_URL_ORDER}/orders`),
        getOne: (id) => http.get(`${BASE_URL_ORDER}/orders/${id}`),
        create: (data) => http.post(`${BASE_URL_ORDER}/orders`, JSON.stringify(data), { headers: { 'Content-Type': 'application/json' } }),
    }
};

export const check200 = (res, name) => check(res, { [`${name} is 200`]: (r) => r.status === 200 });
export const check201 = (res, name) => check(res, { [`${name} is 201`]: (r) => r.status === 201 });
