#!/usr/bin/env bash
set -euo pipefail

# Analyze files in a directory — the hand-written shell equivalent of main.go.
# Each section maps directly to one gloo pipeline; see main.go for the Go form.
#
# Usage: ./analyze-files.sh [directory]   (defaults to ".")

dir=${1:-.}
echo "Analyzing files in: ${dir}" >&2 # gloo: note(status, ...) — diagnostics to stderr

# === File extensions (count by type) ===
# gloo: find(FindType f, FindName '*.*') | While(extension) | Sort() | Uniq(UniqCount)
echo "=== File extensions (count by type) ===" >&2
find "${dir}" -type f -name '*.*' |
  while IFS= read -r file; do
    printf '%s\n' "${file##*.}" # gloo: filepath.Ext then TrimPrefix(".")
  done |
  sort |
  uniq -c

# === Largest files (top 10) ===
# gloo: find(FindType f) | While(sizeAndName) | Sort(Numeric,Reverse,Field 1,Tab) | Head(10)
echo "=== Largest files (top 10) ===" >&2
find "${dir}" -type f |
  while IFS= read -r file; do
    printf '%s\t%s\n' "$(wc -c <"${file}" | tr -d ' ')" "${file}" # gloo: fs.Stat(name).Size()
  done |
  sort -k1,1nr |
  head -10

# === Total size ===
# gloo: find(FindType f) | While(sizeOnly) | Awk(totalSize)
echo "=== Total size ===" >&2
find "${dir}" -type f |
  while IFS= read -r file; do
    wc -c <"${file}" | tr -d ' '
  done |
  awk '{ sum += $1 } END { print "total: " sum " bytes" }'
