package main

import (
	"io"
	"os"
	"strings"
)

// Example_report projects the name/salary columns of a small CSV and totals the
// salary column. Section headers go to io.Discard; only the data lands on stdout.
func Example_report() {
	csv := strings.NewReader(
		"name,dept,salary\n" +
			"alice,eng,120\n" +
			"bob,sales,90\n" +
			"carol,eng,130\n",
	)
	if err := report(csv, os.Stdout, io.Discard); err != nil {
		panic(err)
	}
	// Output:
	// name,salary
	// alice,120
	// bob,90
	// carol,130
	// total: 340
}
