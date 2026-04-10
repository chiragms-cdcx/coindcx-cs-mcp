BINARY          := coindcx-cs-mcp
PKG             := ./cmd/coindcx-cs-mcp
AUDIT_BINARY    := audit-consumer
AUDIT_PKG       := ./cmd/audit-consumer
MIGRATE_BINARY  := audit-migrate
MIGRATE_PKG     := ./cmd/audit-migrate

-include .env
export

.PHONY: build run build-audit-consumer run-audit-consumer migrate tidy clean test test-v help

build:
	go build -o $(BINARY) $(PKG)

run: build
	MCP_HTTP_PORT=8080 MCP_HTTP_HOST=0.0.0.0 ./$(BINARY)

go-run:
	MCP_HTTP_PORT=8080 MCP_HTTP_HOST=0.0.0.0 go run $(PKG)

test:
	go test ./...

test-v:
	go test ./... -v

tidy:
	go mod tidy

clean:
	rm -f $(BINARY) $(AUDIT_BINARY) $(MIGRATE_BINARY)

build-audit-consumer:
	go build -o $(AUDIT_BINARY) $(AUDIT_PKG)

run-audit-consumer: build-audit-consumer
	./$(AUDIT_BINARY)

migrate:
	go run $(MIGRATE_PKG)

help:
	@echo "Targets:"
	@echo "  build              - build $(BINARY)"
	@echo "  run                - build and run server on http://0.0.0.0:8080"
	@echo "  go-run             - run with go run on :8080"
	@echo "  build-audit-consumer - build $(AUDIT_BINARY)"
	@echo "  run-audit-consumer   - build and run audit consumer"
	@echo "  migrate            - run audit schema migration"
	@echo "  test               - run all tests"
	@echo "  test-v             - run all tests (verbose)"
	@echo "  tidy               - go mod tidy"
	@echo "  clean              - remove binaries"
