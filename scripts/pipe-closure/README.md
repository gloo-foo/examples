# Pipe Closure Example

A runnable gloo program that shows how a downstream command which stops reading (`head`, `tail`) tears down the whole pipeline, so upstream generators — even an infinite one like `yes` — stop producing. This is the gloo equivalent of Unix `SIGPIPE`: `seq 1 1000000 | head -n 5` must not generate a million lines.

## How gloo handles it

`head`/`tail` stop pulling once they have their quota; the framework cancels the upstream stream's context, and generators observe the cancellation and stop. No error is reported — early termination is the expected outcome.

## Running

```bash
go run .
```

[`pipe-closure.sh`](pipe-closure.sh) is the equivalent shell, kept side-by-side so each gloo pipeline reads against the shell it mirrors.

## Testing

Every scenario writes to an injected `io.Writer`, so [`main_test.go`](main_test.go) exercises each as a runnable [`Example`](https://pkg.go.dev/testing#hdr-Examples) whose `// Output:` is verified by `go test`:

```bash
go test ./...     # or: make test
```

## Scenarios

Each `head.Head(head.HeadLines(n))` / `tail.Tail(tail.TailLines(n))` is wired with `gloo.Run(source, sink, cmds...)`.

| Scenario | gloo | Shell |
| --- | --- | --- |
| `head` stops after 3 of 10 | `seq.Seq(1, 10)`, `head.Head(head.HeadLines(3))` | `seq 1 10 \| head -n 3` |
| `tail` keeps the last 5 of 100 | `seq.Seq(1, 100)`, `tail.Tail(tail.TailLines(5))` | `seq 1 100 \| tail -n 5` |
| `head` then `grep` | `seq.Seq(1, 20)`, `head.Head(head.HeadLines(5))`, `grep.Grep("3")` | `seq 1 20 \| head -n 5 \| grep 3` |
| infinite `yes`, bounded by `head` | `yes.Yes(yes.YesText("hello"))`, `head.Head(head.HeadLines(3))` | `yes hello \| head -n 3` |
| stacked `head`s | `seq.Seq(1, 100)`, `head.Head(head.HeadLines(50))`, `…(10)`, `…(3)` | `seq 1 100 \| head -n 50 \| head -n 10 \| head -n 3` |

The infinite-`yes` case is the important one: without pipe closure it would never terminate. The stacked-`head` case shows several closure points composing — the smallest quota wins.
