# Aggregate Makefile for the gloo examples. Every immediate subdirectory that
# has its own Makefile is a module group (framework, scripts); `make test` runs
# each group's tests and `make check` runs each group's full quality gate. New
# groups are discovered automatically — no edits needed here.
.DELETE_ON_ERROR:
.DEFAULT_GOAL := help

# Absolute path to this Makefile's directory (trailing slash), so the wildcard
# and the sub-makes resolve from any caller's working directory.
here := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
# Each subdirectory holding a Makefile, as a bare group name.
MODULES := $(patsubst $(here)%/Makefile,%,$(wildcard $(here)*/Makefile))

.PHONY: help
help: ## List available targets
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: $(MODULES) ## Run tests for every module group

.PHONY: $(MODULES)
$(MODULES):
	$(MAKE) -C $(here)$@ test

.PHONY: check
check: $(addprefix check-,$(MODULES)) ## Run the full quality gate for every module group

.PHONY: $(addprefix check-,$(MODULES))
$(addprefix check-,$(MODULES)):
	$(MAKE) -C $(here)$(@:check-%=%) check
