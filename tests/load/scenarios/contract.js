import { Api, check200 } from '../common/api.js';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

export const options = {
    vus: 1,
    iterations: 1,
};

export default function () {
    const res = Api.product.getAll();
    check200(res, 'Products API');

    const products = res.json();
    if (products.length > 0) {
        const p = products[0];
        // Basic "Contract" check: verify essential fields exist
        const hasFields = p.hasOwnProperty('id') && p.hasOwnProperty('sku') && p.hasOwnProperty('price');
        if (!hasFields) {
            console.error('Contract violation: Product missing fields');
        }
    }
}

export function handleSummary(data) {
    return {
        "tests/load/reports/contract_report.html": htmlReport(data),
    };
}
