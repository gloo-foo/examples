# Log Processor Example

A runnable gloo program that scans every `*.log` file in [`logs/`](logs), keeps the error and warning lines, and prints `timestamp,level` (CSV) for each.

It demonstrates the idiomatic gloo replacement for shell's **nested** `while read` loops: the outer loop over files is ordinary Go (`afero.ReadDir`), and each file is processed by its own small gloo pipeline (`grep -iE 'error|warning' | While(timestampLevel)`). The directory is read through an injected `afero.Fs`, so the whole program is testable without touching disk.

## Running

```bash
go run .          # writes the CSV rows to stdout, progress to stderr
```

[`process-logs.sh`](process-logs.sh) is the equivalent hand-written shell version, kept side-by-side so you can read the gloo translation against the shell it replaces.

## Testing

[`main_test.go`](main_test.go) drives [`main.go`](main.go) over an in-memory log directory — a runnable [`Example`](https://pkg.go.dev/testing#hdr-Examples) whose `// Output:` is verified by `go test`:

```bash
go test ./...     # or: make test
```

## Patterns shown

- **Go orchestration, gloo per-file pipelines** — the file loop is plain Go; only the per-file filtering is a gloo pipeline.
- **Extended, case-insensitive grep** — `grep -iE 'error|warning'` is `grep.Grep("error|warning", grep.GrepExtended, grep.GrepIgnoreCase)`. `Grep` matches a fixed string by default; `GrepExtended` is required for the `|` alternation.
- **Injected filesystem** — `afero.ReadDir(fs, dir)` plus `gloo.ByteFileSource(fs, …)` make the run hermetic in tests (`afero.NewMemMapFs()`) and real on disk (`afero.NewOsFs()`).
