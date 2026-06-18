package main

import (
	"os"
	"strings"
)

// Example_frequency ranks a column of access-log IPs by how often each appears
// and prints the top three, biggest first. The output is "<count>\t<value>".
func Example_frequency() {
	log := strings.NewReader(
		"10.0.0.1\n10.0.0.2\n10.0.0.1\n203.0.113.7\n10.0.0.1\n10.0.0.2\n10.0.0.1\n",
	)
	if err := frequency(log, os.Stdout, 3); err != nil {
		panic(err)
	}
	// Output:
	// 4	10.0.0.1
	// 2	10.0.0.2
	// 1	203.0.113.7
}
