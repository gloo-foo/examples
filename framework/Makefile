# Self-documenting quality gate for an illustrative gloo demo module. Run `make`
# (or `make help`).
#
# Every tool is declared in the go.mod `tool` stanza and invoked via `go tool`,
# so there are no global installs and versions are pinned in go.mod/go.sum. These
# modules are documentation/demos: the `check` target proves they stay clean,
# compile, and (for example tests) pass — it does not enforce coverage or build
# release artifacts. No change is complete until `check` exits zero.
.DEFAULT_GOAL := check

GO ?= go
# Production (non-test) Go files — the cognitive-complexity gate runs over these.
SRC := $(shell find . -name '*.go' -not -name '*_test.go' -not -path './vendor/*')

.PHONY: help
help: ## List available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'

## Quality gate

.PHONY: check
check: fmt-check vet lint staticcheck cognit vuln build test ## Full gate: format, vet, lint, staticcheck, complexity<=7, vuln, build, test

.PHONY: fmt
fmt: ## Rewrite all files with the strict formatter (gofumpt)
	$(GO) tool gofumpt -w .

.PHONY: fmt-check
fmt-check: ## Fail if any file is not gofumpt-clean
	@out="$$($(GO) tool gofumpt -l .)"; \
	if [ -n "$$out" ]; then echo "gofumpt would reformat:"; echo "$$out"; exit 1; fi

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: lint
lint: ## Run golangci-lint aggregate analysis
	$(GO) tool golangci-lint run

.PHONY: staticcheck
staticcheck: ## Run staticcheck (zero findings)
	$(GO) tool staticcheck ./...

.PHONY: cognit
cognit: ## Assert cognitive complexity <= 7 for every production function
	@if [ -z "$(SRC)" ]; then exit 0; fi; \
	out="$$($(GO) tool gocognit -over 7 $(SRC))"; \
	if [ -n "$$out" ]; then echo "cognitive complexity > 7:"; echo "$$out"; exit 1; fi

.PHONY: vuln
vuln: ## Scan for known vulnerabilities
	$(GO) tool govulncheck ./...

.PHONY: build
build: ## Compile every package
	$(GO) build ./...

.PHONY: test
test: ## Run tests (including the runnable Example tests)
	$(GO) test ./...

## Utilities

.PHONY: tidy
tidy: ## Tidy and verify module dependencies
	$(GO) mod tidy
	$(GO) mod verify
