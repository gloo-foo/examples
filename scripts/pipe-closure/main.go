// Command pipe-closure demonstrates how downstream commands that stop reading
// (head, tail) tear down the pipeline, so upstream generators — even infinite
// ones like yes — stop producing. It mirrors shell pipe/SIGPIPE behavior.
//
// Each scenario writes to an injected io.Writer, so the demos double as runnable
// Example tests (see main_test.go).
package main

import (
	"fmt"
	"io"
	"os"

	grep "github.com/gloo-foo/cmd-grep"
	head "github.com/gloo-foo/cmd-head"
	seq "github.com/gloo-foo/cmd-seq"
	tail "github.com/gloo-foo/cmd-tail"
	yes "github.com/gloo-foo/cmd-yes"
	gloo "github.com/gloo-foo/framework"
)

func main() {
	scenarios := []func(io.Writer){
		headStops,
		tailKeepsLast,
		headThenGrep,
		yesIsBounded,
		stackedHeads,
	}
	for _, scenario := range scenarios {
		scenario(os.Stdout)
		note(os.Stdout, "\n")
	}
}

// headStops shows head reading only its quota: seq generates ten lines, but the
// pipe closes after head has taken three, so the remaining seven are never sent.
func headStops(w io.Writer) {
	note(w, "=== head stops the pipe after 3 of 10 lines ===\n")
	run(w, seq.Seq(1, 10), head.Head(head.HeadLines(3)))
}

// tailKeepsLast shows tail, which must drain its input to know the end, keeping
// only the final five of a hundred lines.
func tailKeepsLast(w io.Writer) {
	note(w, "=== tail keeps only the last 5 of 100 lines ===\n")
	run(w, seq.Seq(1, 100), tail.Tail(tail.TailLines(5)))
}

// headThenGrep shows head bounding the data a downstream grep ever sees: only
// 1..5 reach grep, of which "3" matches.
func headThenGrep(w io.Writer) {
	note(w, "=== generate 20, head 5, then grep for '3' ===\n")
	run(w, seq.Seq(1, 20), head.Head(head.HeadLines(5)), grep.Grep("3"))
}

// yesIsBounded shows the critical case: yes is infinite, but head closing the
// pipe after three lines stops it — without pipe closure this would never end.
func yesIsBounded(w io.Writer) {
	note(w, "=== yes is infinite; head stops it after 3 lines ===\n")
	run(w, yes.Yes(yes.YesText("hello")), head.Head(head.HeadLines(3)))
}

// stackedHeads shows several heads composing: each closes its own pipe, so the
// smallest quota wins and only three lines survive.
func stackedHeads(w io.Writer) {
	note(w, "=== stacked heads short-circuit to the smallest ===\n")
	run(w, seq.Seq(1, 100),
		head.Head(head.HeadLines(50)),
		head.Head(head.HeadLines(10)),
		head.Head(head.HeadLines(3)))
}

// run wires source → commands → w, exiting on the (here unexpected) error.
func run(w io.Writer, source gloo.Source[[]byte], cmds ...any) {
	if _, err := gloo.Run(source, gloo.ByteWriteTo(w), cmds...); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// note writes a scenario header to w. A failed write to the demo output is not
// actionable here, so the error is deliberately ignored.
func note(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}
