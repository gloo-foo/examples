#!/usr/bin/env bash
set -euo pipefail

# JSON queries — the hand-written shell equivalent of main.go. Section headers go
# to stderr; only the jq results land on stdout.
#
# Usage: ./json-query.sh < data.json

data=$(cat)

echo "=== all names ===" >&2 # gloo: note(status, ...)
# gloo: jq.Jq("-r", ".items[].name")
jq -r '.items[].name' <<<"${data}"

echo "=== engineers ===" >&2
# gloo: jq.Jq("-r", `.items[] | select(.role == "eng") | .name`)
jq -r '.items[] | select(.role == "eng") | .name' <<<"${data}"
