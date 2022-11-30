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

.PHONY: test-leaks
test-leaks:
	@go test -c -o tests
	@for test in $$(go test -list . | grep -E "^(Test|Example)"); do ./tests -test.run "^$$test\$$" &>/dev/null && echo -n "." || echo -e "\n$$test failed"; done


.PHONY: bench
bench: ## run benchmarks
	@go test -v ./... -bench=. -count 2 -benchmem -run=^#

.PHONY: cov
cov: ## run tests and generate coverage report
	@go test -v ./... -coverpkg=./... -race -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: vet
vet: ## check go code
	@go vet ./...

.PHONY: fmt
fmt: ## run gofmt in all project files
	@go fmt ./...

.PHONY: check
check: vet ## check source code
	@staticcheck ./...

.PHONY: deps
deps: ## check dependencies
	@go mod verify

.PHONY: download
download: ## download dependencies
	@go mod download

.PHONY: docs
docs: ## show godocs
	@echo navigate to: http://localhost:6060/pkg/github.com/vitorsalgado/mocha/v3
	@godoc -http=:6060

.PHONY: init
init: tools ## prepare local development environment
ifdef BREW_EXISTS
	brew bundle
endif

.PHONY: tools
tools:
	@echo "tools"
	@make -C tools
