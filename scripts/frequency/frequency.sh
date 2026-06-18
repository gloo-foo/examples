#!/usr/bin/env bash
set -euo pipefail

# Top N most frequent lines — the hand-written shell equivalent of main.go.
# Reads one value per line on stdin and prints the most common, with counts.
#
# Usage: ./frequency.sh [n] < values   (n defaults to 10)

n=${1:-10}

# gloo: sortcmd.Sort() | uniq.Uniq(UniqCount) | <count\tvalue> | sortcmd.Sort(Numeric,Reverse,Field 1) | Head(n)
sort | uniq -c | sort -rn | head -n "${n}"
