# vanity-go Makefile

# Variables
BINARY_NAME := vanity-go
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOVET := $(GOCMD) vet
GOLINT := golangci-lint

# Build parameters
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitHash=$(COMMIT_HASH)"

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## build: Build the binary
.PHONY: build
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

## build-all: Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64
	@GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe
	@echo "Build complete! Binaries are in ./dist/"

## test: Run tests
.PHONY: test
test:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

## test-coverage: Run tests with coverage report
.PHONY: test-coverage
test-coverage: test
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## bench: Run benchmarks
.PHONY: bench
bench:
	$(GOTEST) -bench=. -benchmem ./...

## clean: Clean build files
.PHONY: clean
clean:
	$(GOCLEAN)
	@rm -f $(BINARY_NAME)
	@rm -rf dist/
	@rm -f coverage.out coverage.html

## deps: Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) verify

## tidy: Tidy go.mod
.PHONY: tidy
tidy:
	$(GOMOD) tidy

## fmt: Format code
.PHONY: fmt
fmt:
	$(GOFMT) -s -w .

## vet: Run go vet
.PHONY: vet
vet:
	$(GOVET) ./...

## lint: Run linter
.PHONY: lint
lint:
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## check: Run all checks (fmt, vet, lint, test)
.PHONY: check
check: fmt vet lint test

## run: Run the application
.PHONY: run
run: build
	@if [ -z "$(DOMAIN)" ] || [ -z "$(REPOSITORY)" ]; then \
		echo "Error: DOMAIN and REPOSITORY environment variables are required"; \
		echo "Usage: make run DOMAIN=go.example.com REPOSITORY=https://github.com/username"; \
		exit 1; \
	fi
	./$(BINARY_NAME)

## docker-build: Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t $(BINARY_NAME):latest .
	docker build -t $(BINARY_NAME):$(VERSION) .

## docker-run: Run Docker container
.PHONY: docker-run
docker-run:
	@if [ -z "$(DOMAIN)" ] || [ -z "$(REPOSITORY)" ]; then \
		echo "Error: DOMAIN and REPOSITORY environment variables are required"; \
		echo "Usage: make docker-run DOMAIN=go.example.com REPOSITORY=https://github.com/username"; \
		exit 1; \
	fi
	docker run -p 8080:8080 -e DOMAIN=$(DOMAIN) -e REPOSITORY=$(REPOSITORY) $(BINARY_NAME):latest

## install: Install the binary
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to $(GOPATH)/bin/"
	@cp $(BINARY_NAME) $(GOPATH)/bin/

## uninstall: Uninstall the binary
.PHONY: uninstall
uninstall:
	@echo "Removing $(BINARY_NAME) from $(GOPATH)/bin/"
	@rm -f $(GOPATH)/bin/$(BINARY_NAME)

## release: Create a new release (requires TAG variable)
.PHONY: release
release:
	@if [ -z "$(TAG)" ]; then \
		echo "Error: TAG variable is required"; \
		echo "Usage: make release TAG=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating release $(TAG)..."
	git tag -a $(TAG) -m "Release $(TAG)"
	git push origin $(TAG)

## serve: Run with example configuration
.PHONY: serve
serve:
	DOMAIN=localhost:8080 REPOSITORY=https://github.com/example $(GOCMD) run main.go

## watch: Run with file watching (requires entr)
.PHONY: watch
watch:
	@if command -v entr >/dev/null 2>&1; then \
		find . -name '*.go' | entr -r make run; \
	else \
		echo "entr not installed. Install it to use watch mode."; \
		exit 1; \
	fi

.PHONY: all
all: clean deps fmt vet lint test build