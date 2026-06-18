// Command log-processor scans every *.log file in a directory, keeps the
// error/warning lines, and prints "timestamp,level" for each. It shows the
// idiomatic replacement for shell's nested "while read" loops: orchestrate the
// per-file pipelines in Go, each built from gloo commands. The directory is read
// through an injected afero.Fs so the whole program is testable in memory
// (see main_test.go).
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	grep "github.com/gloo-foo/cmd-grep"
	while "github.com/gloo-foo/cmd-while"
	gloo "github.com/gloo-foo/framework"
	"github.com/spf13/afero"
)

func main() {
	if err := process(afero.NewOsFs(), os.Stdout, os.Stderr, "logs"); err != nil {
		fmt.Fprintf(os.Stderr, "log-processor: %v\n", err)
		os.Exit(1)
	}
}

// process runs the per-file pipeline over every *.log file in dir, writing the
// CSV rows to out and progress to status.
func process(fs afero.Fs, out, status io.Writer, dir string) error {
	entries, err := afero.ReadDir(fs, dir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := processEntry(fs, out, status, dir, entry); err != nil {
			return err
		}
	}
	return nil
}

// processEntry filters one directory entry: non-.log files (and subdirectories)
// are skipped; a log file is grepped for error/warning lines, each rewritten as
// "timestamp,level".
// shell: grep -iE 'error|warning' file | while read; do echo "$1,$2"; done
func processEntry(fs afero.Fs, out, status io.Writer, dir string, entry os.FileInfo) error {
	if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".log") {
		return nil
	}
	path := filepath.Join(dir, entry.Name())
	note(status, "Processing %s\n", path)

	source := gloo.ByteFileSource(fs, []gloo.File{gloo.File(path)})
	_, err := gloo.Run(
		source, gloo.ByteWriteTo(out),
		grep.Grep("error|warning", grep.GrepExtended, grep.GrepIgnoreCase),
		while.While(timestampLevel),
	)
	return err
}

// note writes a diagnostic line to w. A failed write to a progress/status
// stream is not actionable here, so the error is deliberately ignored.
func note(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

// timestampLevel turns a log line into "<field1>,<field2>" (timestamp,level),
// skipping lines with fewer than two whitespace-separated fields.
func timestampLevel(line []byte) ([]byte, error) {
	fields := strings.Fields(string(line))
	if len(fields) < 2 {
		return nil, nil
	}
	return fmt.Appendf(nil, "%s,%s", fields[0], fields[1]), nil
}
