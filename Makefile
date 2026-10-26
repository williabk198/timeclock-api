# Makefile for the Time Clock API
#
# Run `make help` (or just `make`) to list available targets.
#
# Running a microservice (e.g. `make run-admin`) requires environment
# variables to be set, matching what's configured in .vscode/launch.json
# and .zed/debug.json:
#   DB_URI           postgres connection string (all services)
#   PORT             HTTP port to listen on (all services)
#   GRPC_PORT        gRPC port to listen on (auth service only)
#   AUTH_GRPC_ADDR   host:port of the auth service's gRPC listener (admin service only)

MODULE   := github.com/williabk198/timeclock
BIN_DIR  := bin
MICROS   := admin auth

# NOTE: build-%/run-% are intentionally not listed here. Listing a stem-expanded
# name (e.g. "build-admin") as a literal .PHONY prerequisite creates an explicit,
# recipe-less target that shadows the build-%/run-% pattern rules below, causing
# `make build-admin` to report "Nothing to be done". They have no same-named file
# to produce, so make always treats them as out-of-date anyway.
.PHONY: all help build test vet fmt generate tidy tools clean

.DEFAULT_GOAL := help

all: build ## Build all microservice binaries

help: ## Show this help message
	@grep -E '^[a-zA-Z0-9_.%-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

build: $(addprefix build-,$(MICROS)) ## Build all microservice binaries into bin/

build-%: ## Build a single microservice binary, e.g. `make build-admin`
	go build -o $(BIN_DIR)/$* ./cmd/micros/$*

run-%: ## Run a single microservice without building a binary, e.g. `make run-admin`
	go run ./cmd/micros/$*

test: ## Run the full test suite
	go test ./...

vet: ## Run go vet across the module
	go vet ./...

fmt: ## Format all Go source files
	gofmt -l -w .

generate: ## Regenerate generated code (e.g. protobuf/gRPC stubs) from source
	go generate ./...

tidy: ## Tidy go.mod/go.sum
	go mod tidy

tools: ## Install the protoc-gen-go/protoc-gen-go-grpc plugins used by `make generate`
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)
