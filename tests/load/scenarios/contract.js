import { check, group } from 'k6';
import { Api } from '../common/api.js';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

// Load Swagger specifications
const productSwagger = JSON.parse(open('../../../product-service/docs/swagger.json'));
const orderSwagger = JSON.parse(open('../../../order-service/docs/swagger.json'));

export const options = {
    vus: 1,
    iterations: 1,
};

/**
 * Basic Schema Validator for k6
 */
function validateSchema(data, schemaName, swagger) {
    // Find the schema in definitions. Swag sometimes uses full paths.
    let schema = swagger.definitions[schemaName];
    if (!schema) {
        // Try finding by suffix if not found exactly (e.g. "domain.Product" in full path)
        const key = Object.keys(swagger.definitions).find(k => k.endsWith(schemaName));
        if (key) {
            schema = swagger.definitions[key];
        }
    }

    if (!schema) {
        console.error(`Schema ${schemaName} not found in Swagger definitions`);
        return false;
    }

    // If it's just a basic type (like an enum or special object), it might not have properties
    if (!schema.properties) {
        if (schema.type === 'object') return typeof data === 'object';
        if (schema.type === 'string') return typeof data === 'string';
        return true; 
    }

    const properties = schema.properties;
    const required = schema.required || [];
    let isValid = true;

    // Check required fields
    required.forEach(field => {
        if (!data.hasOwnProperty(field)) {
            console.error(`Contract Violation: Missing required field '${field}' for schema '${schemaName}'`);
            isValid = false;
        }
    });

    // Check types
    for (const [prop, details] of Object.entries(properties)) {
        if (data.hasOwnProperty(prop) && data[prop] !== null) {
            const actualValue = data[prop];
            const actualType = typeof actualValue;
            let expectedType = details.type;

            if (expectedType === 'integer') expectedType = 'number';
            
            if (details['$ref']) {
                const nestedSchemaName = details['$ref'].split('/').pop();
                isValid = validateSchema(actualValue, nestedSchemaName, swagger) && isValid;
            } else if (actualType !== expectedType && expectedType !== undefined) {
                 console.error(`Contract Violation: Field '${prop}' expected type '${expectedType}', got '${actualType}'`);
                 isValid = false;
            }
        }
    }

    return isValid;
}

export default function () {
    group('Product Service Contract', () => {
        const res = Api.product.getAll();
        check(res, { 'product list status is 200': (r) => r.status === 200 });

        const products = res.json();
        if (products.length > 0) {
            check(products[0], {
                'product matches defined schema': (p) => validateSchema(p, 'domain.Product', productSwagger)
            });
        }
    });

    group('Order Service Contract', () => {
        const res = Api.order.getAll();
        check(res, { 'order list status is 200': (r) => r.status === 200 });

        const orders = res.json();
        if (orders.length > 0) {
            check(orders[0], {
                'order matches defined schema': (o) => validateSchema(o, 'domain.Order', orderSwagger)
            });
        }
    });
}

export function handleSummary(data) {
    return {
        "tests/load/reports/contract_report.html": htmlReport(data),
    };
}
