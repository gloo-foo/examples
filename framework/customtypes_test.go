package gloo_test

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	gloo "github.com/gloo-foo/framework"
	"github.com/gloo-foo/framework/patterns"
)

// These examples show how to build commands around custom, strongly-typed
// parameters using gloo.NewParameters[P, F]: P is the positional type (instead
// of gloo.File) and F is the flags struct. The commands themselves are ordinary
// gloo.Command[[]byte, []byte] values built from the canonical patterns.

// runOn executes a byte command over input and returns its joined output.
func runOn(cmd gloo.Command[[]byte, []byte], input string) string {
	ctx := context.Background()
	source := gloo.ByteReaderSource([]io.Reader{strings.NewReader(input)})
	lines, err := gloo.Collect(ctx, cmd.Execute(ctx, source.Stream(ctx)))
	if err != nil {
		return fmt.Sprintf("error: %v\n", err)
	}
	var b strings.Builder
	for _, line := range lines {
		b.Write(line)
		b.WriteByte('\n')
	}
	return b.String()
}

// ============================================================================
// Custom struct parameters + a flag, parsed via ...any
// ============================================================================

// SearchPattern is a strongly-typed search string.
type SearchPattern string

// FilterConfig is a custom struct used as a positional parameter.
type FilterConfig struct {
	Field    string
	MinValue int
}

// myFlags holds the command's flags.
type myFlags struct{ CaseSensitive bool }

// CaseSensitivity is a sealed bool flag.
type CaseSensitivity bool

const (
	CaseSensitive   CaseSensitivity = true
	CaseInsensitive CaseSensitivity = false
)

func (c CaseSensitivity) Configure(f *myFlags) { f.CaseSensitive = bool(c) }

// StronglyTypedCommand keeps lines containing pattern and tags each with the
// first matching FilterConfig's field. Positional FilterConfig values and the
// CaseSensitivity flag are parsed from the variadic parameters.
func StronglyTypedCommand(pattern SearchPattern, parameters ...any) gloo.Command[[]byte, []byte] {
	p := gloo.NewParameters[FilterConfig, myFlags](parameters...)
	configs := p.Typed
	caseSensitive := p.Flags.CaseSensitive

	return patterns.Expand(func(line []byte) ([][]byte, error) {
		s := string(line)
		if !contains(string(pattern), s, caseSensitive) {
			return nil, nil
		}
		return tagFirstMatch(configs, s), nil
	})
}

// contains reports whether s contains pattern, honoring case sensitivity.
func contains(pattern, s string, caseSensitive bool) bool {
	if caseSensitive {
		return strings.Contains(s, pattern)
	}
	return strings.Contains(strings.ToLower(s), strings.ToLower(pattern))
}

// tagFirstMatch returns s tagged with the field of the first config that
// applies — one whose MinValue is unset, or whose MinValue appears in s — and no
// lines when none apply.
func tagFirstMatch(configs []FilterConfig, s string) [][]byte {
	for _, c := range configs {
		if c.MinValue > 0 && !strings.Contains(s, strconv.Itoa(c.MinValue)) {
			continue
		}
		return [][]byte{fmt.Appendf(nil, "[%s] %s", c.Field, s)}
	}
	return nil
}

func Example_stronglyTypedCommand() {
	cmd := StronglyTypedCommand(
		SearchPattern("error"),
		FilterConfig{Field: "log", MinValue: 100},
		FilterConfig{Field: "alert", MinValue: 200},
		CaseInsensitive,
	)
	fmt.Print(runOn(cmd, "ERROR 100: Something went wrong\nINFO: all good\nERROR 200: Critical issue\n"))
	// Output:
	// [log] ERROR 100: Something went wrong
	// [alert] ERROR 200: Critical issue
}

// ============================================================================
// Explicit type safety at the signature (no ...any)
// ============================================================================

// ExplicitlyTypedCommand is like StronglyTypedCommand but takes its configs as a
// typed variadic, so the compiler rejects anything that is not a FilterConfig.
func ExplicitlyTypedCommand(pattern SearchPattern, configs ...FilterConfig) gloo.Command[[]byte, []byte] {
	return patterns.Expand(func(line []byte) ([][]byte, error) {
		s := string(line)
		if !strings.Contains(s, string(pattern)) {
			return nil, nil
		}
		return tagFirstMatch(configs, s), nil
	})
}

func Example_explicitlyTypedCommand() {
	cmd := ExplicitlyTypedCommand(
		SearchPattern("test"),
		FilterConfig{Field: "field1", MinValue: 10},
		FilterConfig{Field: "field2", MinValue: 20},
	)
	fmt.Print(runOn(cmd, "test line with 10\nanother test with 20\n"))
	// Output:
	// [field1] test line with 10
	// [field2] another test with 20
}

// ============================================================================
// A custom positional type as a Source (like gloo.File, but a URL)
// ============================================================================

// URL is a custom positional type.
type URL string

// FetchURLs is a Source that emits one "Fetching: <url>" line per URL argument.
func FetchURLs(urls ...any) gloo.Source[[]byte] {
	p := gloo.NewParameters[URL, struct{}](urls...)
	return &urlSource{urls: p.Typed}
}

type urlSource struct{ urls []URL }

func (s *urlSource) Stream(ctx context.Context) gloo.Stream[[]byte] {
	return gloo.Generate(ctx, func(_ context.Context, send func([]byte) bool, _ func(error)) {
		for _, u := range s.urls {
			if !send([]byte("Fetching: " + string(u))) {
				return
			}
		}
	})
}

func Example_customPositionalType() {
	src := FetchURLs(
		URL("https://example.com"),
		URL("https://github.com"),
		URL("https://google.com"),
	)
	ctx := context.Background()
	lines, _ := gloo.Collect(ctx, src.Stream(ctx))
	for _, line := range lines {
		fmt.Println(string(line))
	}
	// Output:
	// Fetching: https://example.com
	// Fetching: https://github.com
	// Fetching: https://google.com
}
