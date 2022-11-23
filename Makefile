MOCHA_NO_COLOR := "0"
BREW_EXISTS := $(shell brew --version 2>/dev/null)

.ONESHELL:
.DEFAULT_GOAL := help

# allow user specific optional overrides
-include Makefile.overrides

export

.PHONY: help
help: ## show help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: test
test: ## run tests
	@go test -race -v ./... -race

.PHONY: bench
bench: ## run benchmarks
	@go test -v ./... -bench=. -count 2 -benchmem -run=^#

.PHONY: cov
cov: ## run tests and generate coverage report
	@go test -v ./... -coverpkg=./... -race -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

vet: ## check go code
	@go vet ./...

fmt: ## run gofmt in all project files
	@go fmt ./...

check: vet ## check source code
	@staticcheck ./...

deps: ## check dependencies
	@go mod verify

download: ## download dependencies
	@go mod download

.PHONY: docs
docs: ## show godocs
	@echo navigate to: http://localhost:6060/pkg/github.com/vitorsalgado/mocha/v3
	@godoc -http=:6060

init: ## prepare local development environment
	@echo "tools"
	@make -C tools
ifdef YOUR_PROGRAM_VERSION
	brew bundle
endif
