package main

import (
	"io"
	"os"
	"strings"
)

// Example_query lists every item's name, then just the engineers. Section
// headers go to io.Discard; only the jq results land on stdout.
func Example_query() {
	doc := strings.NewReader(`{
  "items": [
    {"name": "alice", "role": "eng"},
    {"name": "bob",   "role": "sales"},
    {"name": "carol", "role": "eng"}
  ]
}`)
	if err := query(doc, os.Stdout, io.Discard); err != nil {
		panic(err)
	}
	// Output:
	// alice
	// bob
	// carol
	// alice
	// carol
}
