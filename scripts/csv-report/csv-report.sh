#!/usr/bin/env bash
set -euo pipefail

# CSV report — the hand-written shell equivalent of main.go. Section headers go
# to stderr; only the data lands on stdout.
#
# Usage: ./csv-report.sh < data.csv

data=$(cat)

echo "=== name, salary ===" >&2 # gloo: note(status, ...)
# gloo: cut.Cut(CutFields(1, 3), CutDelimiter(","))
cut -d, -f1,3 <<<"${data}"

echo "=== total payroll ===" >&2
# gloo: cut.Cut(CutFields(3), CutDelimiter(",")) | awk(sumColumn)
cut -d, -f3 <<<"${data}" | awk '{ sum += $1 } END { print "total: " sum }'
