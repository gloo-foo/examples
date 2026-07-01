# Self-documenting quality gate for an illustrative gloo demo command. Run `make`
# (or `make help`).
#
# Tools are resolved from ${GOBIN} only -- the pinned set installed by
# nicerobot/tools.build (go-tooling/tools.txt), the same mechanism the shared
# library Makefile uses. There is no per-module tool stanza or vendoring. This
# module is a runnable documentation demo (package main): the `check` target
# proves it stays clean and compiles — it does not enforce coverage or build
# release artifacts. No change is complete until `check` exits zero.
.DEFAULT_GOAL := check

GO    ?= go
GOBIN := $(strip $(shell $(GO) env GOBIN))
gobin-or-die = $(or $(GOBIN),$(error GOBIN is not set -- install the pinned tool set via nicerobot/tools.build (go-tooling/install-tools.sh); the gate resolves tools from ${GOBIN} only, no PATH/GOPATH fallback))
GOFUMPT       = $(gobin-or-die)/gofumpt
GOLANGCI_LINT = $(gobin-or-die)/golangci-lint
STATICCHECK   = $(gobin-or-die)/staticcheck
GOCOGNIT      = $(gobin-or-die)/gocognit
GOVULNCHECK   = $(gobin-or-die)/govulncheck
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
	$(GOFUMPT) -w .

.PHONY: fmt-check
fmt-check: ## Fail if any file is not gofumpt-clean
	@out="$$($(GOFUMPT) -l .)"; \
	if [ -n "$$out" ]; then echo "gofumpt would reformat:"; echo "$$out"; exit 1; fi

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: lint
lint: ## Run golangci-lint aggregate analysis
	$(GOLANGCI_LINT) run

.PHONY: staticcheck
staticcheck: ## Run staticcheck (zero findings)
	$(STATICCHECK) ./...

.PHONY: cognit
cognit: ## Assert cognitive complexity <= 7 for every production function
	@if [ -z "$(SRC)" ]; then exit 0; fi; \
	out="$$($(GOCOGNIT) -over 7 $(SRC))"; \
	if [ -n "$$out" ]; then echo "cognitive complexity > 7:"; echo "$$out"; exit 1; fi

.PHONY: vuln
vuln: ## Scan for known vulnerabilities
	$(GOVULNCHECK) ./...

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
