#!/usr/bin/env bash
set -euo pipefail

# Process log files — the hand-written shell equivalent of main.go.
# Emits "timestamp,level" CSV rows to stdout; progress goes to stderr.
#
# Usage: ./process-logs.sh [directory]   (defaults to "logs")

dir=${1:-logs}

# gloo: afero.ReadDir(fs, dir) then a per-file pipeline for each *.log
for file in "${dir}"/*.log; do
  [[ -e "${file}" ]] || continue
  echo "Processing ${file}" >&2 # gloo: fmt.Fprintf(status, "Processing %s\n", path)

  # gloo: grep.Grep("error|warning", GrepExtended, GrepIgnoreCase) | While(timestampLevel)
  grep -iE 'error|warning' "${file}" \
    | awk '{ print $1 "," $2 }'
done
