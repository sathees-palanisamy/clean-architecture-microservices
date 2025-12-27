import { sleep } from 'k6';
import { Api, check200, check201 } from '../common/api.js';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

export const options = {
    stages: [
        { duration: '2m', target: 50 },  // normal load
        { duration: '3m', target: 100 }, // stress point
        { duration: '2m', target: 200 }, // breaking point
        { duration: '2m', target: 0 },   // recovery
    ],
    thresholds: {
        http_req_duration: ['p(95)<1000'],
        http_req_failed: ['rate<0.05'],
    },
};

export default function () {
    const productsRes = Api.product.getAll();
    const products = productsRes.json();

    if (products && products.length > 0) {
        const product = products[Math.floor(Math.random() * products.length)];
        Api.product.getOne(product.id);
        
        const orderData = {
            user_id: Math.floor(Math.random() * 1000),
            product_id: product.id,
            quantity: 1
        };
        Api.order.create(orderData);
    }
    sleep(1);
}

export function handleSummary(data) {
    return {
        "tests/load/reports/stress_report.html": htmlReport(data),
    };
}
