.PHONY: build run up down clean check-quality

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
