BINARY_NAME := inscribe
BUILD_DIR := bin
GO := go

.PHONY: build test test-coverage lint clean run

build:
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/inscribe

test:
	$(GO) test ./... -v

test-coverage:
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)
