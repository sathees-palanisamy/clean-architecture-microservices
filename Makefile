.PHONY: build run up down clean check-quality test test-verbose generate test-integration load-test smoke-test stress-test soak-test spike-test contract-test benchmark-test fault-test

build:
	go build -o bin/product-service ./product-service/cmd/product-service/main.go
	go build -o bin/order-service ./order-service/cmd/order-service/main.go

run-product:
	DB_HOST=localhost DB_PORT=5433 SERVER_PORT=8081 go run ./product-service/cmd/product-service/main.go

run-order:
	DB_HOST=localhost DB_PORT=5434 SERVER_PORT=8082 PRODUCT_SERVICE_URL=http://localhost:8081 go run ./order-service/cmd/order-service/main.go

up:
	docker-compose up --build -d

down:
	docker-compose down

clean:
	rm -rf bin
	docker-compose down -v

# Helper to check if everything compiles
check:
	go work sync
	cd pkg && go vet ./...
	cd product-service && go vet ./...
	cd order-service && go vet ./...

test:
	go test ./pkg/...
	go test ./product-service/...
	go test ./order-service/...

test-verbose:
	go test -v ./pkg/...
	go test -v ./product-service/...
	go test -v ./order-service/...

generate:
	cd pkg && go generate ./...
	cd product-service && go generate ./...
	cd order-service && go generate ./...

test-integration:
	cd order-service/internal/integration_test && go test -v .
	cd product-service/internal/integration_test && go test -v .

load-test:
	k6 run tests/load/load_test.js

smoke-test:
	k6 run tests/load/scenarios/smoke.js

stress-test:
	k6 run tests/load/scenarios/stress.js

soak-test:
	k6 run tests/load/scenarios/soak.js

spike-test:
	k6 run tests/load/scenarios/spike.js

contract-test:
	k6 run tests/load/scenarios/contract.js

benchmark-test:
	k6 run tests/load/scenarios/benchmark.js

fault-test:
	k6 run tests/load/scenarios/fault_injection.js
