# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.work and mod files
COPY go.work go.work.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY product-service/go.mod product-service/go.sum ./product-service/
COPY order-service/go.mod order-service/go.sum ./order-service/

# Download dependencies
RUN go mod download

# Copy source code
COPY pkg ./pkg
COPY order-service ./order-service

# Build the application
WORKDIR /app/order-service
RUN go build -o /app/order-service-bin ./cmd/order-service/main.go

# Final stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/order-service-bin .
COPY --from=builder /app/order-service/schema ./schema

EXPOSE 8082
CMD ["./order-service-bin"]
