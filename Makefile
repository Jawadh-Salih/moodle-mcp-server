BINARY     := moodle-mcp
CMD        := ./cmd/moodle-mcp
GO         := go
GOFLAGS    := -trimpath
LDFLAGS    := -s -w

.PHONY: all build test bench bench-api bench-tools coverage clean lint help

all: build

## build: compile the binary
build:
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD)

## test: run all unit tests
test:
	$(GO) test -race ./...

## test-v: run all unit tests with verbose output
test-v:
	$(GO) test -race -v ./...

## bench: run all benchmarks
bench:
	$(GO) test -bench=. -benchmem ./...

## bench-api: run API layer benchmarks only
bench-api:
	$(GO) test -bench=. -benchmem ./internal/api/...

## bench-tools: run tool handler benchmarks only
bench-tools:
	$(GO) test -bench=. -benchmem ./internal/tools/...

## coverage: generate HTML coverage report
coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report written to coverage.html"

## coverage-text: print coverage summary to terminal
coverage-text:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out

## lint: run go vet and staticcheck (install staticcheck separately)
lint:
	$(GO) vet ./...
	@which staticcheck > /dev/null && staticcheck ./... || echo "staticcheck not installed — run: go install honnef.co/go/tools/cmd/staticcheck@latest"

## tidy: tidy and verify go modules
tidy:
	$(GO) mod tidy
	$(GO) mod verify

## clean: remove build artefacts
clean:
	rm -f $(BINARY) coverage.out coverage.html

## help: show this help
help:
	@grep -E '^## ' Makefile | sed 's/## //'
