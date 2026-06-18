# Frequency Example

A runnable gloo program for the classic "top talkers" question: **which values occur most often?** Feed it one value per line — IP addresses from an access log, error codes, URLs, status lines — and it prints the `n` most frequent, with counts, biggest first.

This is the one-liner everyone reaches for at the shell:

```bash
sort | uniq -c | sort -rn | head
```

…expressed as a single gloo composition.

## Running

```bash
cut -d' ' -f1 access.log | go run . 5      # top 5 client IPs
go run . < values.txt                      # top 10 (default)
```

[`frequency.sh`](frequency.sh) is the equivalent hand-written shell version, kept side-by-side so you can read the gloo translation against the shell it replaces.

## Testing

```bash
go test ./...     # or: make test
```

[`main_test.go`](main_test.go) ranks a small column of IPs as a runnable [`Example`](https://pkg.go.dev/testing#hdr-Examples).

## Patterns shown

- **`sort | uniq -c`** — `sortcmd.Sort()` then `uniq.Uniq(uniq.UniqCount)` collapse equal lines into counts (`uniq` only merges *adjacent* lines, so the sort comes first).
- **Numeric ranking on a field** — `uniq -c` emits a right-justified `"   <count> <value>"`; a tiny `awk` step rewrites it to `"<count>\t<value>"` so `sortcmd.Sort(SortNumeric, SortReverse, SortField(1), SortDelimiterTab)` keys cleanly on the count even when the value contains spaces. (`SortNumeric` parses the whole key, so the field split is what makes the numeric sort correct.)
- **Bounded output** — `head.Head(head.HeadLines(n))` keeps only the top `n`.

> Values with equal counts are reported in an unspecified order.
