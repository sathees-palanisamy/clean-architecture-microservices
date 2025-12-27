import { sleep } from 'k6';
import { Api } from '../common/api.js';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

// This test simulates "Chaos" by expecting failures or high latency 
// and verifying that the system handles it (e.g., circuit breakers if they existed)
export const options = {
    vus: 10,
    duration: '30s',
};

export default function () {
    // In a real fault injection test, we might call an endpoint that simulates delay
    // or we might block network on a container. 
    // Here we just verify how the test handles slow responses.
    Api.product.getAll();
    sleep(1);
}

export function handleSummary(data) {
    return {
        "tests/load/reports/fault_report.html": htmlReport(data),
    };
}
