// Command csv-report reads a CSV with a header row and prints two summaries of
// it: a projection of chosen columns, and the sum of a numeric column. It shows
// how gloo's field-aware commands (cut) compose with a custom awk program for
// everyday tabular wrangling — the kind of thing you'd otherwise reach for
// `cut`/`awk` at the shell to do.
//
// Usage: csv-report < data.csv
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"

	awk "github.com/gloo-foo/cmd-awk"
	cut "github.com/gloo-foo/cmd-cut"
	gloo "github.com/gloo-foo/framework"
)

func main() {
	if err := report(os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "csv-report: %v\n", err)
		os.Exit(1)
	}
}

// report reads the whole CSV once, then runs each summary over its own source.
func report(in io.Reader, out, status io.Writer) error {
	data, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	if err := projectColumns(data, out, status); err != nil {
		return err
	}
	return totalSalary(data, out, status)
}

// projectColumns prints just the name and salary columns (1 and 3).
// shell: cut -d, -f1,3
func projectColumns(data []byte, out, status io.Writer) error {
	note(status, "=== name, salary ===\n")
	return run(data, out, cut.Cut(cut.CutFields(1, 3), cut.CutDelimiter(",")))
}

// totalSalary sums the salary column. The header cell ("salary") is not numeric,
// so the awk program skips it without any special-casing.
// shell: cut -d, -f3 | awk '{s+=$1} END{print s}'
func totalSalary(data []byte, out, status io.Writer) error {
	note(status, "=== total payroll ===\n")
	return run(
		data, out,
		cut.Cut(cut.CutFields(3), cut.CutDelimiter(",")),
		awk.Awk(&sumColumn{}),
	)
}

// run builds a fresh source from data and wires it through cmds to out.
func run(data []byte, out io.Writer, cmds ...any) error {
	src := gloo.ByteReaderSource([]io.Reader{bytes.NewReader(data)})
	_, err := gloo.Run(src, gloo.ByteWriteTo(out), cmds...)
	return err
}

// note writes a diagnostic header to w; a failed write is not actionable here.
func note(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

// sumColumn totals the first field of every line that parses as an integer,
// emitting "total: <sum>" at the end.
type sumColumn struct {
	awk.SimpleProgram
	sum int64
}

func (p *sumColumn) Action(ctx *awk.Context) (string, bool) {
	if n, err := strconv.ParseInt(ctx.Field(1), 10, 64); err == nil {
		p.sum += n
	}
	return "", false
}

func (p *sumColumn) End(*awk.Context) (string, error) {
	return fmt.Sprintf("total: %d", p.sum), nil
}
