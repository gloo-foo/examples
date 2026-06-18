package pipeline_test

import (
	"context"
	"fmt"
	"os"
	"strconv"

	. "github.com/gloo-foo/cmd-awk"
	. "github.com/gloo-foo/cmd-cat"
	. "github.com/gloo-foo/cmd-sort"
	gloo "github.com/gloo-foo/framework"
	"github.com/spf13/afero"
)

// sumProgram sums the second field of each line (awk '{sum+=$2} END{...}').
type sumProgram struct {
	SimpleProgram
	sum int
}

func (p *sumProgram) Action(ctx *Context) (string, bool) {
	if n, err := strconv.Atoi(ctx.Field(2)); err == nil {
		p.sum += n
	}
	return "", false
}

func (p *sumProgram) End(*Context) (string, error) {
	return fmt.Sprintf("total: %d", p.sum), nil
}

// extractSecondField prints the second field of each line (awk '{print $2}').
type extractSecondField struct{ SimpleProgram }

func (extractSecondField) Action(ctx *Context) (string, bool) {
	return ctx.Field(2), true
}

// fileSource reads a file as a byte stream — the source for these pipelines.
func fileSource(path string) gloo.Source[[]byte] {
	return gloo.ByteFileSource(afero.NewOsFs(), []gloo.File{gloo.File(path)})
}

// run wires source → commands → stdout.
func run(src gloo.Source[[]byte], cmds ...any) {
	if _, err := gloo.Run(src, gloo.ByteWriteTo(os.Stdout), cmds...); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}

// ExamplePipe_catSortAwk composes the stages with gloo.Pipe, then runs them.
func ExamplePipe_catSortAwk() {
	pipeline := gloo.Pipe[[]byte, []byte, []byte](
		gloo.Pipe[[]byte, []byte, []byte](Cat(), Sort()),
		Awk(&sumProgram{}),
	)
	run(fileSource("testdata/fruits.txt"), pipeline)
	// Output:
	// total: 36
}

// ExampleChain_catSortAwk builds the same pipeline with the fluent Chain API.
func ExampleChain_catSortAwk() {
	_, err := gloo.Chain(fileSource("testdata/fruits.txt")).
		To(Cat()).To(Sort()).To(Awk(&sumProgram{})).
		Sink(gloo.ByteWriteTo(os.Stdout))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	// Output:
	// total: 36
}

// ExampleChain_sortAwkExtract sorts the file, then prints each second field.
func ExampleChain_sortAwkExtract() {
	run(fileSource("testdata/fruits.txt"), Sort(), Awk(extractSecondField{}))
	// Output:
	// 3
	// 6
	// 7
	// 2
	// 5
	// 1
	// 8
	// 4
}

// ExampleChain_catNumberSort numbers the lines, then sorts the numbered output.
func ExampleChain_catNumberSort() {
	run(fileSource("testdata/fruits.txt"), Cat(CatNumberLines), Sort())
	// Output:
	//      1	banana 5
	//      2	apple 3
	//      3	cherry 8
	//      4	banana 2
	//      5	apple 7
	//      6	cherry 1
	//      7	date 4
	//      8	apple 6
}

// ExampleChain_empty shows empty input yields empty output.
func ExampleChain_empty() {
	ctx := context.Background()
	src := gloo.SliceSource([][]byte(nil))
	pipeline := gloo.Pipe[[]byte, []byte, []byte](
		gloo.Pipe[[]byte, []byte, []byte](Cat(), Sort()),
		Awk(extractSecondField{}),
	)
	lines, _ := gloo.Collect(ctx, pipeline.Execute(ctx, src.Stream(ctx)))
	fmt.Printf("lines: %d\n", len(lines))
	// Output:
	// lines: 0
}
