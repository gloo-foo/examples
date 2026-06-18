// Command json-query runs jq filters over a JSON document on stdin: it extracts
// a field from every element of an array, then filters and projects with a
// select(). It shows that gloo composes structured-data tools as readily as the
// classic line-oriented ones — the gloo form of piping through `jq`.
//
// Usage: json-query < data.json
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	jq "github.com/gloo-foo/cmd-jq"
	gloo "github.com/gloo-foo/framework"
)

func main() {
	if err := query(os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "json-query: %v\n", err)
		os.Exit(1)
	}
}

// query reads the whole JSON document once, then applies each jq filter to it.
func query(in io.Reader, out, status io.Writer) error {
	data, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	if err := run(data, out, status, "=== all names ===", ".items[].name"); err != nil {
		return err
	}
	return run(data, out, status, "=== engineers ===", `.items[] | select(.role == "eng") | .name`)
}

// run applies one jq filter (raw output) to data, writing results to out and the
// section header to status.
func run(data []byte, out, status io.Writer, header, filter string) error {
	note(status, "%s\n", header)
	src := gloo.ByteReaderSource([]io.Reader{bytes.NewReader(data)})
	_, err := gloo.Run(src, gloo.ByteWriteTo(out), jq.Jq("-r", filter))
	return err
}

// note writes a diagnostic header to w; a failed write is not actionable here.
func note(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}
