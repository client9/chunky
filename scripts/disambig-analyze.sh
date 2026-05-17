#!/usr/bin/env bash
set -euo pipefail
word="${1:?usage: disambig-analyze.sh <word>}"
corpus="data/*.pos.jsonl"

echo "=== next-tag patterns for '$word' ==="
go run ./cmd/lexrules -feat word+nexttag -sort-count $corpus | grep -E "^\s+[0-9]+\. ${word}\+"

echo ""
echo "=== prev-tag patterns for '$word' ==="
go run ./cmd/lexrules -feat word+prevtag -sort-count $corpus | grep -E "^\s+[0-9]+\. ${word}\+"
