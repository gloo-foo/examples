# File Statistics Example

A runnable gloo program that analyzes a directory and prints three statistics:

1. **File extensions (count by type)** — groups files by extension and counts them (`sort | uniq -c`).
2. **Largest files (top 10)** — the ten biggest files, largest first.
3. **Total size** — the combined size of every file, summed with a custom `awk` program.

It demonstrates a more complex composition than the single-command `cmd-*/examples`: each section is its own gloo pipeline (`find → While → sort/uniq/head/awk`), and the per-file transform uses native Go (`afero.Fs.Stat`, `filepath.Ext`) inside the `While` body instead of shell's `ls`/`awk` field-parsing.

## Running

```bash
go run . [directory]   # defaults to the current directory
```

[`analyze-files.sh`](analyze-files.sh) is the equivalent hand-written shell pipeline, kept side-by-side so you can read the gloo translation against the shell it replaces.

## Testing

The analysis in [`main.go`](main.go) is written against an injected `afero.Fs`, so [`main_test.go`](main_test.go) exercises it hermetically over an in-memory filesystem with known file sizes — a runnable [`Example`](https://pkg.go.dev/testing#hdr-Examples) whose `// Output:` is verified by `go test`:

```bash
go test ./...    # or: make test
```

## Patterns shown

- **Injected filesystem** — `find.Find(dir, find.FindFs(fs), …)` and `fs.Stat` make the whole program testable without touching disk.
- **Field-aware numeric sort** — `sort -k1,1nr` is `sortcmd.Sort(SortNumeric, SortReverse, SortField(1), SortDelimiterTab)`: the size is field 1 of a `"<size>\t<name>"` line, so the numeric comparison applies to the size and not the trailing path.
- **Custom `awk` program** — `totalSize` accumulates field 1 across all lines and emits the total in `End`.
