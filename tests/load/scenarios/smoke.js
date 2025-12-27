import { sleep } from 'k6';
import { Api, check200 } from '../common/api.js';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

export const options = {
    vus: 1,
    duration: '10s',
    thresholds: {
        http_req_failed: ['rate<0.01'],
    },
};

export default function () {
    check200(Api.product.health(), 'Product Health');
    check200(Api.order.health(), 'Order Health');
    sleep(1);
}

export function handleSummary(data) {
    return {
        "tests/load/reports/smoke_report.html": htmlReport(data),
    };
}
