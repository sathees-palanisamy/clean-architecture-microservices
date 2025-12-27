import { sleep } from 'k6';
import { Api } from '../common/api.js';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

export const options = {
    stages: [
        { duration: '10s', target: 20 },  // base load
        { duration: '20s', target: 500 }, // sudden spike
        { duration: '1m', target: 500 },  // sustain spike
        { duration: '20s', target: 20 },  // recovery
        { duration: '10s', target: 0 },
    ],
    thresholds: {
        http_req_failed: ['rate<0.1'],
    },
};

export default function () {
    Api.product.health();
    sleep(0.5);
}

export function handleSummary(data) {
    return {
        "tests/load/reports/spike_report.html": htmlReport(data),
    };
}
