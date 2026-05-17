# unix-socket-path-limit — Makefile
#
# Reader validation flow:
#   make test              fast unit tests, no Docker
#   make test-integration  qemu in a real container (needs Docker)
#   make verify            both, in order — the full proof

GO                  ?= go
INTEGRATION_TIMEOUT ?= 15m

.DEFAULT_GOAL := help
.PHONY: help test test-integration verify vet tidy check-docker

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN{FS=":.*?## "}{printf "  %-18s %s\n", $$1, $$2}'

test: ## Run the fast table-driven unit tests (no Docker)
	@$(GO) test -v ./runtimedir/...

test-integration: check-docker ## Run the qemu-in-Testcontainers tests (needs Docker)
	@$(GO) test -v -count=1 -timeout $(INTEGRATION_TIMEOUT) -tags integration ./integration/...

verify: test test-integration ## Full validation: unit tests, then integration

vet: ## go vet across all packages, including the integration build tag
	@$(GO) vet -tags integration ./...

tidy: ## Sync go.mod / go.sum
	@$(GO) mod tidy

check-docker: ## Fail early with a clear message if Docker is unreachable
	@docker info >/dev/null 2>&1 || { \
		echo "Docker is not running — start it and retry."; exit 1; }
