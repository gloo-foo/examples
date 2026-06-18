# Pipeline Demo Example

Runnable [`Example`](https://pkg.go.dev/testing#hdr-Examples) tests that compose `cmd-*` stages two different ways and verify the output with `go test`. Unlike the other `scripts/*` examples there is no `main.go`: the demos live entirely in [`pipeline_test.go`](pipeline_test.go), each with a verified `// Output:`.

## Running

```bash
go test ./...     # or: make test
go test -v        # see each example's name and output
```

## What it shows

Both APIs build the same `cat | sort | awk` pipeline over [`testdata/fruits.txt`](testdata/fruits.txt):

- **`gloo.Pipe`** — explicit, type-parameterized composition: `gloo.Pipe[[]byte, []byte, []byte](gloo.Pipe(Cat(), Sort()), Awk(&sumProgram{}))`.
- **`gloo.Chain`** — the fluent equivalent: `gloo.Chain(src).To(Cat()).To(Sort()).To(Awk(&sumProgram{})).Sink(...)`.

It also demonstrates custom `awk` programs (`sumProgram` sums a field; `extractSecondField` projects one) and that an empty input yields empty output.
