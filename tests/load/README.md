# Advanced k6 Testing Suite

This directory contains a suite of performance and reliability tests for the microservices.

## Test Types

| Test Type | Objective | Command |
|-----------|-----------|---------|
| **Smoke** | Verify basic connectivity and health | `make smoke-test` |
| **Stress** | Find system breaking points | `make stress-test` |
| **Soak** | Check for memory leaks over time | `make soak-test`|
| **Spike** | Test resilience to sudden traffic bursts | `make spike-test` |
| **Contract** | Validate API response structure | `make contract-test` |
| **Benchmark** | Establish performance baselines | `make benchmark-test` |
| **Fault Injection** | Simulate degraded conditions | `make fault-test` |

## Reports

Every test execution generates an HTML report in the `reports/` directory.
You can open these in your browser to view detailed metrics and visualizations.

## Configuration

The tests use the following environment variables (defined in your shell or `Makefile`):

- `PRODUCT_SERVICE_URL`: Defaults to `http://localhost:8081`
- `ORDER_SERVICE_URL`: Defaults to `http://localhost:8082`

## Prerequisites

- [k6](https://k6.io/docs/getting-started/installation/) installed.
- Microservices must be running (`make run` or `docker-compose up`).
