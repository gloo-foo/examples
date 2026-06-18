#!/usr/bin/env bash
# Note: no `pipefail` here on purpose — these pipelines rely on the upstream
# generator being killed by SIGPIPE when head/tail stops reading, which pipefail
# would (correctly, but unhelpfully here) report as a pipeline failure.
set -eu

# Demonstrates pipe closure — the shell equivalent of main.go.
# A downstream command (head/tail) that stops reading closes the pipe, so
# upstream generators receive SIGPIPE and stop. Each block mirrors one gloo
# scenario; see main.go.

echo "=== head stops the pipe after 3 of 10 lines ==="
seq 1 10 | head -n 3 # gloo: seq.Seq(1, 10), head.Head(head.HeadLines(3))
echo

echo "=== tail keeps only the last 5 of 100 lines ==="
seq 1 100 | tail -n 5 # gloo: seq.Seq(1, 100), tail.Tail(tail.TailLines(5))
echo

echo "=== generate 20, head 5, then grep for '3' ==="
seq 1 20 | head -n 5 | grep 3 # gloo: seq.Seq(1, 20), head.Head(head.HeadLines(5)), grep.Grep("3")
echo

echo "=== yes is infinite; head stops it after 3 lines ==="
yes hello | head -n 3 # gloo: yes.Yes(yes.YesText("hello")), head.Head(head.HeadLines(3))
echo

echo "=== stacked heads short-circuit to the smallest ==="
seq 1 100 | head -n 50 | head -n 10 | head -n 3 # gloo: three head.Head(head.HeadLines(...))
echo
