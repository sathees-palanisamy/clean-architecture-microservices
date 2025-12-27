# Observability Guide for Ops Engineers

This guide outlines how to monitor and troubleshoot the Order and Product microservices using the implemented OpenTelemetry (OTEL) and Grafana stack.

## 1. Stack Overview
- **Instrumentation**: Go OTEL SDK (HTTP, DB, Usecases).
- **Collector**: OpenTelemetry Collector (receives OTLP, exports to Prometheus/Tempo).
- **Metrics**: Prometheus (HTTP latency, error rates, DB pool stats).
- **Tracing**: Grafana Tempo (Distributed tracing).
- **Logging**: Grafana Loki (Log-Trace correlation).
- **Visualization**: Grafana.

## 2. Monitoring Dashboards
Access the **Microservices Observability** dashboard at `http://localhost:3000`.

### Key Metrics to Watch:
- **HTTP Latency (P95)**: Should be < 500ms. If it spikes, check Tempo for slow spans.
- **Error Rate**: Should be < 1%. Spikes usually indicate downstream issues or DB failures.
- **DB Connection Pool**: If `InUse` reaches `MaxOpenConns` (25), requests will queue and latency will increase.

## 3. Distributed Tracing (Tempo)
When investigating a slow request:
1. Find the `trace_id` in the application logs or the "Explore" tab.
2. Search for the trace in Tempo.
3. Identify which span (HTTP, Usecase, or SQL) is taking the most time.

## 4. Alerting
The following alerts are pre-configured:
- **High Latency**: Triggered if P95 > 500ms for 5 minutes.
- **Error Spike**: Triggered if error rate > 5% for 2 minutes.
- **DB Saturation**: Triggered if connection pool usage > 90% for 1 minute.

## 5. Kubernetes Deployment
For production, use the provided sidecar manifest in `k8s/otel-collector-sidecar.yaml`.
Ensure the following environment variables are set per service:
- `OTEL_SERVICE_NAME`: The name of the service (e.g., `order-service`).
- `OTEL_EXPORTER_OTLP_ENDPOINT`: The endpoint of the collector (e.g., `localhost:4317` for sidecar).

## 6. Local Development
Run the entire stack using docker-compose:
```bash
docker-compose up -d
```
All ports are exposed locally:
- Grafana: 3000
- Prometheus: 9090
- Tempo: 3200
- OTEL Collector: 4317 (gRPC) / 4318 (HTTP)
