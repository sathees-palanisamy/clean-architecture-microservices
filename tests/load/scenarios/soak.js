import { sleep } from 'k6';
import { Api } from '../common/api.js';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

export const options = {
    stages: [
        { duration: '1m', target: 30 },  // ramp up
        { duration: '2m', target: 30 },  // sustain load for long time
        { duration: '1m', target: 0 },   // ramp down
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'],
        http_req_failed: ['rate<0.01'],
    },
};

export default function () {
    Api.product.getAll();
    Api.order.getAll();
    sleep(2);
}

export function handleSummary(data) {
    return {
        "tests/load/reports/soak_report.html": htmlReport(data),
    };
}
