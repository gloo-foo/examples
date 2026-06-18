package main

import (
	"io"
	"os"

	"github.com/spf13/afero"
)

// Example_analyze runs the full analysis over an in-memory filesystem with
// known file sizes, so the output is deterministic and the example doubles as a
// test. Section headers go to io.Discard; only the data lands on stdout.
func Example_analyze() {
	fs := afero.NewMemMapFs()
	write := func(name, content string) {
		if err := afero.WriteFile(fs, name, []byte(content), 0o644); err != nil {
			panic(err)
		}
	}
	write("data/alpha.txt", "aaaa")            // 4 bytes
	write("data/beta.go", "bb")                // 2 bytes
	write("data/gamma.go", "gggggggg")         // 8 bytes
	write("data/delta.md", "dddddddddddddddd") // 16 bytes

	if err := analyze(fs, os.Stdout, io.Discard, "data"); err != nil {
		panic(err)
	}
	// Output:
	//       2 go
	//       1 md
	//       1 txt
	// 16	data/delta.md
	// 8	data/gamma.go
	// 4	data/alpha.txt
	// 2	data/beta.go
	// total: 30 bytes
}
