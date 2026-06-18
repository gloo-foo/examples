package main

import "os"

// The scenarios are deterministic, so each is a runnable Example whose
// // Output: is verified by `go test`.

func Example_headStops() {
	headStops(os.Stdout)
	// Output:
	// === head stops the pipe after 3 of 10 lines ===
	// 1
	// 2
	// 3
}

func Example_tailKeepsLast() {
	tailKeepsLast(os.Stdout)
	// Output:
	// === tail keeps only the last 5 of 100 lines ===
	// 96
	// 97
	// 98
	// 99
	// 100
}

func Example_headThenGrep() {
	headThenGrep(os.Stdout)
	// Output:
	// === generate 20, head 5, then grep for '3' ===
	// 3
}

func Example_yesIsBounded() {
	yesIsBounded(os.Stdout)
	// Output:
	// === yes is infinite; head stops it after 3 lines ===
	// hello
	// hello
	// hello
}

func Example_stackedHeads() {
	stackedHeads(os.Stdout)
	// Output:
	// === stacked heads short-circuit to the smallest ===
	// 1
	// 2
	// 3
}
