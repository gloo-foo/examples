// Command file-stats analyzes the files in a directory using gloo pipelines:
// find → transform-per-file (While) → sort/uniq/head/awk. It shows how native
// Go (afero.Fs.Stat, filepath) replaces shell's ls/awk field-parsing inside
// While bodies, and how the analysis is written against an injected filesystem
// so it can be exercised hermetically (see main_test.go).
//
// Usage: file-stats [dir]   (defaults to ".")
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	awk "github.com/gloo-foo/cmd-awk"
	find "github.com/gloo-foo/cmd-find"
	head "github.com/gloo-foo/cmd-head"
	sortcmd "github.com/gloo-foo/cmd-sort"
	uniq "github.com/gloo-foo/cmd-uniq"
	while "github.com/gloo-foo/cmd-while"
	gloo "github.com/gloo-foo/framework"
	"github.com/spf13/afero"
)

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	if err := analyze(afero.NewOsFs(), os.Stdout, os.Stderr, dir); err != nil {
		fmt.Fprintf(os.Stderr, "file-stats: %v\n", err)
		os.Exit(1)
	}
}

// analyze runs every statistics section over dir, writing results to out and
// human-readable section headers to status. The filesystem is injected so the
// same code runs against the real OS filesystem (main) and an in-memory one
// (the example test).
func analyze(fs afero.Fs, out, status io.Writer, dir string) error {
	note(status, "Analyzing files in: %s\n", dir)
	sections := []func(afero.Fs, io.Writer, io.Writer, string) error{
		extensionCounts,
		largestFiles,
		totalBytes,
	}
	for _, section := range sections {
		if err := section(fs, out, status, dir); err != nil {
			return err
		}
	}
	return nil
}

// extensionCounts counts files per extension, sorted by extension.
// shell: find -name '*.*' | while read; do echo ext; done | sort | uniq -c
func extensionCounts(fs afero.Fs, out, status io.Writer, dir string) error {
	note(status, "\n=== File extensions (count by type) ===\n")
	return run(
		out,
		find.Find(dir, find.FindFs(fs), find.FindType("f"), find.FindName("*.*")),
		while.While(extension),
		sortcmd.Sort(),
		uniq.Uniq(uniq.UniqCount),
	)
}

// largestFiles lists the ten largest files, biggest first.
// shell: find -type f | while read; do echo "size\tname"; done | sort -k1,1nr | head -10
func largestFiles(fs afero.Fs, out, status io.Writer, dir string) error {
	note(status, "\n=== Largest files (top 10) ===\n")
	return run(
		out,
		find.Find(dir, find.FindFs(fs), find.FindType("f")),
		while.While(sizeAndName(fs, status)),
		sortcmd.Sort(sortcmd.SortNumeric, sortcmd.SortReverse, sortcmd.SortField(1), sortcmd.SortDelimiterTab),
		head.Head(head.HeadLines(10)),
	)
}

// totalBytes sums the sizes of every file.
// shell: find -type f | while read; do echo size; done | awk '{s+=$1} END{...}'
func totalBytes(fs afero.Fs, out, status io.Writer, dir string) error {
	note(status, "\n=== Total size ===\n")
	return run(
		out,
		find.Find(dir, find.FindFs(fs), find.FindType("f")),
		while.While(sizeOnly(fs, status)),
		awk.Awk(&totalSize{}),
	)
}

// run wires source → commands → out.
func run(out io.Writer, source gloo.Source[[]byte], cmds ...any) error {
	_, err := gloo.Run(source, gloo.ByteWriteTo(out), cmds...)
	return err
}

// note writes a diagnostic line to w. A failed write to a progress/status
// stream is not actionable here, so the error is deliberately ignored.
func note(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

// extension maps a filename to its extension without the dot (empty if none).
func extension(line []byte) ([]byte, error) {
	return []byte(strings.TrimPrefix(filepath.Ext(string(line)), ".")), nil
}

// sizeAndName maps a filename to "<size>\t<name>" via fs.Stat. A file that
// cannot be stat'd (e.g. one that vanished between find and stat) is logged to
// status and skipped by emitting no line, so one bad file never aborts the run.
func sizeAndName(fs afero.Fs, status io.Writer) func([]byte) ([]byte, error) {
	return func(line []byte) ([]byte, error) {
		name := string(line)
		info, err := fs.Stat(name)
		if err != nil {
			note(status, "skipping %s: %v\n", name, err)
			return nil, nil
		}
		return fmt.Appendf(nil, "%d\t%s", info.Size(), name), nil
	}
}

// sizeOnly maps a filename to just its size in bytes, skipping unreadable files.
func sizeOnly(fs afero.Fs, status io.Writer) func([]byte) ([]byte, error) {
	return func(line []byte) ([]byte, error) {
		info, err := fs.Stat(string(line))
		if err != nil {
			note(status, "skipping %s: %v\n", string(line), err)
			return nil, nil
		}
		return fmt.Appendf(nil, "%d", info.Size()), nil
	}
}

// totalSize sums the first field of every line (the sizes).
type totalSize struct {
	awk.SimpleProgram
	sum int64
}

func (p *totalSize) Action(ctx *awk.Context) (string, bool) {
	if size, err := strconv.ParseInt(ctx.Field(1), 10, 64); err == nil {
		p.sum += size
	}
	return "", false
}

func (p *totalSize) End(*awk.Context) (string, error) {
	return fmt.Sprintf("total: %d bytes", p.sum), nil
}
