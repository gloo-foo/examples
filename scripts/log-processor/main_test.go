package main

import (
	"io"
	"os"

	"github.com/spf13/afero"
)

// Example_process runs the processor over an in-memory logs directory, so the
// CSV output is deterministic. afero.ReadDir returns entries sorted by name, so
// app.log is processed before system.log. Progress goes to io.Discard; only the
// "timestamp,level" rows land on stdout.
func Example_process() {
	fs := afero.NewMemMapFs()
	write := func(name, content string) {
		if err := afero.WriteFile(fs, name, []byte(content), 0o644); err != nil {
			panic(err)
		}
	}
	write("logs/app.log", "2024-01-01 ERROR Database connection failed\n"+
		"2024-01-01 INFO Application started\n"+
		"2024-01-01 WARNING Low disk space\n")
	write("logs/system.log", "2024-01-01 INFO System healthy\n"+
		"2024-01-01 ERROR Memory exhausted\n"+
		"2024-01-01 WARNING CPU usage high\n")
	write("logs/notes.txt", "this file is ignored\n")

	if err := process(fs, os.Stdout, io.Discard, "logs"); err != nil {
		panic(err)
	}
	// Output:
	// 2024-01-01,ERROR
	// 2024-01-01,WARNING
	// 2024-01-01,ERROR
	// 2024-01-01,WARNING
}
