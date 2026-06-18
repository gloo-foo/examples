# JSON Query Example

A runnable gloo program that runs [`jq`](https://jqlang.github.io/jq/) filters over a JSON document on stdin — extracting a field from every element of an array, then filtering and projecting with a `select()`. It shows that gloo composes structured-data tools as readily as the classic line-oriented ones.

```bash
jq -r '.items[].name' data.json
jq -r '.items[] | select(.role == "eng") | .name' data.json
```

## Running

```bash
go run . < data.json
```

[`json-query.sh`](json-query.sh) is the equivalent hand-written shell version, kept side-by-side so you can read the gloo translation against the shell it replaces.

## Testing

```bash
go test ./...     # or: make test
```

[`main_test.go`](main_test.go) applies both filters to a small multi-line JSON document as a runnable [`Example`](https://pkg.go.dev/testing#hdr-Examples).

## Patterns shown

- **Structured queries** — `jq.Jq("-r", ".items[].name")` is a full jq filter; `-r` emits raw strings (no surrounding quotes). The framework hands the whole document to jq, so multi-line JSON works as-is.
- **Filter + project** — `jq.Jq("-r", `.items[] | select(.role == "eng") | .name`)` keeps only matching objects and projects one field, all in the jq expression.
- **Reusing one input for several queries** — the document is read once into memory, and each query runs over its own fresh `gloo.ByteReaderSource`.
