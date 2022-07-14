MOCHA_COLOR := 1

.ONESHELL:
.DEFAULT_GOAL := help

export

.PHONY: help basic

help: ## show help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

basic: ## run basic example
	@go run basic/main.go
