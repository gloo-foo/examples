# CSV Report Example

A runnable gloo program for everyday tabular wrangling: read a CSV with a header row, **project a couple of columns**, and **sum a numeric column**. It is the gloo form of the `cut`/`awk` one-liners you'd otherwise type at the shell.

```bash
cut -d, -f1,3 data.csv          # name, salary
cut -d, -f3 data.csv | awk '{s+=$1} END{print s}'   # total payroll
```

## Running

```bash
go run . < data.csv
```

[`csv-report.sh`](csv-report.sh) is the equivalent hand-written shell version, kept side-by-side so you can read the gloo translation against the shell it replaces.

## Testing

```bash
go test ./...     # or: make test
```

[`main_test.go`](main_test.go) runs both summaries over a small in-memory CSV as a runnable [`Example`](https://pkg.go.dev/testing#hdr-Examples).

## Patterns shown

- **Field projection** — `cut.Cut(cut.CutFields(1, 3), cut.CutDelimiter(","))` selects columns 1 and 3 (1-based), the direct equivalent of `cut -d, -f1,3`.
- **Custom `awk` accumulator** — `sumColumn` totals the first field across all lines in `Action` and emits the total in `End`. The header cell (`salary`) is not numeric, so `strconv.ParseInt` fails and it is skipped — no special-casing needed.
- **Reusing one input for several pipelines** — the CSV is read once into memory, and each summary runs over its own fresh `gloo.ByteReaderSource`.
