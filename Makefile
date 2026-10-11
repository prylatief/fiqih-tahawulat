.PHONY: build run dev test lint fmt clean sync build-linux build-mcp run-mcp run-mcp-http bench

APP_NAME = fiqih-tahawulat
BIN_DIR = bin

build:
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/api

run:
	go run ./cmd/api

dev:
	air

test:
	go test -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run

fmt:
	go fmt ./... && go mod tidy

clean:
	rm -rf $(BIN_DIR)/ coverage.out

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 ./cmd/api

build-mcp:
	go build -o $(BIN_DIR)/$(APP_NAME)-mcp ./cmd/mcp

run-mcp:
	go run ./cmd/mcp --transport stdio

run-mcp-http:
	go run ./cmd/mcp --transport http --addr :3001

bench:
	go test -bench=. -benchmem ./internal/repository/
