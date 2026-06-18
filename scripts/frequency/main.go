// Command frequency prints the N most frequent lines in its input — the classic
// "top talkers" pipeline (sort | uniq -c | sort -rn | head) as a single gloo
// composition. Feed it one value per line (IP addresses, error codes, URLs,
// status lines) and it reports the most common, with counts, biggest first.
//
// Usage: frequency [n] < values   (n defaults to 10)
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	awk "github.com/gloo-foo/cmd-awk"
	head "github.com/gloo-foo/cmd-head"
	sortcmd "github.com/gloo-foo/cmd-sort"
	uniq "github.com/gloo-foo/cmd-uniq"
	gloo "github.com/gloo-foo/framework"
)

func main() {
	n, err := topN(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "frequency: %v\n", err)
		os.Exit(1)
	}
	if err := frequency(os.Stdin, os.Stdout, n); err != nil {
		fmt.Fprintf(os.Stderr, "frequency: %v\n", err)
		os.Exit(1)
	}
}

// topN reads the optional count argument, defaulting to 10.
func topN(args []string) (int, error) {
	if len(args) == 0 {
		return 10, nil
	}
	return strconv.Atoi(args[0])
}

// frequency ranks the lines of in by descending frequency and writes the top n
// as "<count>\t<value>" to out:
//
//	sort | uniq -c | <to "count\tvalue"> | sort -k1,1nr | head -n
func frequency(in io.Reader, out io.Writer, n int) error {
	src := gloo.ByteReaderSource([]io.Reader{in})
	_, err := gloo.Run(
		src, gloo.ByteWriteTo(out),
		sortcmd.Sort(),
		uniq.Uniq(uniq.UniqCount),
		awk.Awk(countTab{}),
		sortcmd.Sort(sortcmd.SortNumeric, sortcmd.SortReverse, sortcmd.SortStableSort, sortcmd.SortField(1), sortcmd.SortDelimiterTab),
		head.Head(head.HeadLines(n)),
	)
	return err
}

// countTab rewrites uniq -c's right-justified "   <count> <value>" into
// "<count>\t<value>", so the following numeric sort keys cleanly on the count
// field even when the value itself contains spaces.
type countTab struct{ awk.SimpleProgram }

func (countTab) Action(c *awk.Context) (string, bool) {
	value := make([]string, 0, c.NF)
	for i := 2; i <= c.NF; i++ {
		value = append(value, c.Field(i))
	}
	return c.Field(1) + "\t" + strings.Join(value, " "), true
}
