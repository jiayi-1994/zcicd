.PHONY: all build clean test lint generate migrate

SERVICES := auth project workflow deploy quality artifact system cron
BUILD_DIR := ./bin
GO := go

all: build

build:
	@for svc in $(SERVICES); do \
		echo "Building $$svc..."; \
		$(GO) build -o $(BUILD_DIR)/$$svc ./cmd/$$svc/; \
	done

build-%:
	$(GO) build -o $(BUILD_DIR)/$* ./cmd/$*/

clean:
	rm -rf $(BUILD_DIR)

test:
	$(GO) test ./... -v -cover

lint:
	golangci-lint run ./...

generate:
	$(GO) generate ./...
	wire ./cmd/...

migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

swagger:
	swag init -g cmd/auth/main.go -o api/openapi --parseDependency

proto:
	protoc --go_out=. --go-grpc_out=. api/proto/**/*.proto

docker-build:
	@for svc in $(SERVICES); do \
		docker build -t zcicd-$$svc:latest -f docker/Dockerfile.server --build-arg SERVICE=$$svc .; \
	done
