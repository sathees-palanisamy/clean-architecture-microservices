import { sleep } from 'k6';
import { Api } from '../common/api.js';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

export const options = {
    vus: 10,
    duration: '30s',
    thresholds: {
        http_req_duration: ['p(95)<200'], // Strict threshold for benchmarking
    },
};

export default function () {
    Api.product.health();
    sleep(1);
}

export function handleSummary(data) {
    return {
        "tests/load/reports/benchmark_report.html": htmlReport(data),
    };
}
