# Gloo Framework Examples

Runnable examples for **building commands** with the [gloo framework](https://github.com/gloo-foo/framework), focused on custom, strongly-typed parameters. These are [`Example`](https://pkg.go.dev/testing#hdr-Examples) tests: each has a verified `// Output:`, so `go test` proves they stay correct and they render as runnable examples in godoc.

## Running

```bash
go test ./...                                   # run every example (or: make test)
go test -v -run Example_stronglyTypedCommand    # run one example, verbosely
```

## Examples ([customtypes_test.go](customtypes_test.go))

- `Example_stronglyTypedCommand` — a command built around a custom positional struct (`FilterConfig`) and a sealed flag (`CaseSensitivity`), parsed from variadic parameters via `gloo.NewParameters[P, F]`.
- `Example_explicitlyTypedCommand` — the same idea with a typed variadic (`...FilterConfig`), so the compiler rejects anything that is not a `FilterConfig`.
- `Example_customPositionalType` — a custom positional type (`URL`) used to drive a `gloo.Source`, the way `gloo.File` drives file sources.

## Pattern

Each example:

1. Defines the strongly-typed parameter/flag types.
2. Builds an ordinary `gloo.Command[[]byte, []byte]` from the canonical `patterns` helpers (e.g. `patterns.Expand`).
3. Runs it over in-memory input with `gloo.Collect` / `gloo.ByteReaderSource` and asserts the output.

## Documentation

Full framework documentation lives in the framework repository:

- [gloo framework](https://github.com/gloo-foo/framework) — overview, quick start, and API reference.
