# examples

[![CI](https://github.com/gloo-foo/examples/actions/workflows/ci.yml/badge.svg)](https://github.com/gloo-foo/examples/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Runnable demonstrations of the [gloo framework](https://github.com/gloo-foo/framework) — richer compositions than the single-command examples in each `cmd-*` repository. Every example is verified by `go test` (the `scripts/*` programs via in-memory [`Example`](https://pkg.go.dev/testing#hdr-Examples) tests), so the demos cannot silently drift from the API.

Each example is its own Go module (its own `go.mod`), so they build and test independently.

## Examples

- [`framework/`](framework) — building commands with the framework: custom strongly-typed parameters, sealed flags, and custom positional source types.
- [`scripts/`](scripts) — script-like programs that compose `cmd-*` commands into multi-stage pipelines:
  - [`file-stats`](scripts/file-stats) — directory statistics (`find → While → sort/uniq/head/awk`) over an injected filesystem.
  - [`log-processor`](scripts/log-processor) — extract error/warning rows from many log files (Go orchestration + per-file `grep | While`).
  - [`pipe-closure`](scripts/pipe-closure) — how `head`/`tail` tear down a pipeline and stop upstream generators (gloo's `SIGPIPE`).
  - [`pipeline-demo`](scripts/pipeline-demo) — composing stages with `gloo.Pipe` and the fluent `gloo.Chain` API.
  - [`frequency`](scripts/frequency) — the "top talkers" one-liner: rank the most frequent lines (`sort | uniq -c | sort -rn | head`).
  - [`csv-report`](scripts/csv-report) — tabular wrangling: project columns and sum a numeric column (`cut` + `awk`).
  - [`json-query`](scripts/json-query) — extract and filter structured data with `jq`.

## Running

```bash
make test     # run every example module's tests
make check    # run each module's full quality gate (format, vet, lint, vuln, build, test)
make          # list available targets
```

Each module can also be run on its own, e.g. `cd scripts/file-stats && go run .` or `make test`.
