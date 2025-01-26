sinclude .env

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## install/tools: installs all dev tools
.PHONY: install/tools
install/tools:
	@echo "Installing dev tools"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/favadi/protoc-go-inject-tag@latest
	@go install github.com/yoheimuta/protolint/cmd/protolint@latest
	@echo "Successfully installed dev tools"

## gen/pb: generates protobuf files
.PHONY: gen/pb
gen/pb:
	@echo Starting generate pb
	@mkdir -p ./internal/clients
	@protoc \
	  --proto_path=../ingestor/pkg/api/aggregator \
	  --go_out=internal/clients \
	  --go-grpc_out=internal/clients \
	  ../ingestor/pkg/api/aggregator/aggregator.proto
	@echo Successfully generated proto

	@echo Starting inject tags
	@find ./internal/clients -name "*.pb.go" -exec protoc-go-inject-tag -input="{}" \;
	@echo Successfully injected tags

## test: runs all tests
.PHONY: test
test:
	@go test -v -race ./...

## test/coverage: runs all tests and opens coverage report
.PHONY: test/coverage
test/coverage:
	@mkdir -p ./.coverage
	@go test -v -race ./... -coverprofile=./.coverage/coverage.out -coverpkg=./... -covermode=atomic

## migrate/create: creates a new migration file, takes name as argument
.PHONY: migrate/create
migrate/create:
	@echo "Creating migration file"
	@goose -dir=./database/migrations -s create $(or $(name),"file") sql

## migrate/up: runs all migrations
.PHONY: migrate/up
migrate/up:
	@echo "Running migrations"
	@goose -dir=./database/migrations postgres $(DB_WRITE_DSN) up

## migrate/down: rolls back all migrations
.PHONY: migrate/down
migrate/down:
	@echo "Rolling back migrations"
	@goose -dir=./database/migrations postgres $(DB_WRITE_DSN) down

## seed/up: runs all seeders
.PHONY: seed/up
seed/up:
	@echo "Running seeders"
	@goose -dir=./database/seed -no-versioning postgres $(DB_WRITE_DSN) up

## seed/down: rolls back all seeders
.PHONY: seed/down
seed/down:
	@echo "Rolling back seeders"
	@goose -dir=./database/seed -no-versioning postgres $(DB_WRITE_DSN) down

## lint/go: runs go linters
.PHONY: lint/go
lint/go:
	@echo "Running go linters"
	@golangci-lint run ./... --fix

## lint/proto: runs protobuf linters
.PHONY: lint/proto
lint/proto:
	@echo "Running protobuf linters"
	@protolint lint ./internal/clients/*/*.proto

## run: runs the server
.PHONY: run
run:
	@trap 'kill 0' EXIT; \
	go run ./cmd/server || true
