#!/usr/bin/env python3
"""
Export NLTK corpora to JSONL format compatible with tag-pos.py.

Each output line is {"title": "...", "extract": "..."} — same format
as fetch-wiki-dump.py so tag-pos.py can process it unchanged.

Supported corpora:
  gutenberg  — 18 classic literary texts (Austen, Melville, Carroll, etc.)
  reuters    — ~10k Reuters newswire articles (journalistic register)
  all        — both of the above, in separate output files

Usage:
  python3 scripts/fetch-nltk-corpus.py --corpus gutenberg
  python3 scripts/fetch-nltk-corpus.py --corpus reuters
  python3 scripts/fetch-nltk-corpus.py --corpus all

Output files land in data/:
  data/gutenberg.jsonl
  data/reuters.jsonl
"""

import argparse
import json
import re
import sys
from pathlib import Path

import nltk


CHUNK_WORDS = 500  # max words per record; keeps spaCy memory reasonable


def chunk_text(text, max_words=CHUNK_WORDS):
    """Split text into chunks of at most max_words words, breaking on whitespace."""
    words = text.split()
    for i in range(0, len(words), max_words):
        yield " ".join(words[i:i + max_words])


def fetch_gutenberg(out_path):
    nltk.download("gutenberg", quiet=True)
    from nltk.corpus import gutenberg

    print(f"Writing {out_path} ...", file=sys.stderr)
    count = 0
    with open(out_path, "w") as f:
        for fileid in gutenberg.fileids():
            # Derive a readable title from the filename (e.g. "austen-emma.txt")
            name = fileid.replace(".txt", "").replace("-", " ").title()
            raw = gutenberg.raw(fileid)
            # Strip leading title/header lines (first non-empty paragraph is usually
            # a byline or title that spaCy would mis-segment)
            raw = re.sub(r"^\s*\[.*?\]\s*", "", raw, flags=re.DOTALL)
            raw = raw.strip()
            for chunk in chunk_text(raw):
                f.write(json.dumps({"title": name, "extract": chunk}) + "\n")
                count += 1
    print(f"  {len(gutenberg.fileids())} books → {count} chunks → {out_path}", file=sys.stderr)


def fetch_reuters(out_path):
    nltk.download("reuters", quiet=True)
    from nltk.corpus import reuters

    print(f"Writing {out_path} ...", file=sys.stderr)
    # Use test+train splits; deduplicate by fileid
    fileids = reuters.fileids()
    count = 0
    with open(out_path, "w") as f:
        for fileid in fileids:
            raw = reuters.raw(fileid).strip()
            if not raw:
                continue
            # Reuters raw text starts with the title in ALL CAPS on its own line.
            # Pull it out as the title field.
            lines = raw.splitlines()
            title = lines[0].strip().title() if lines else fileid
            body = "\n".join(lines[1:]).strip() if len(lines) > 1 else raw
            if not body:
                body = raw
            for chunk in chunk_text(body):
                f.write(json.dumps({"title": title, "extract": chunk}) + "\n")
                count += 1
    print(f"  {len(fileids)} articles → {count} chunks → {out_path}", file=sys.stderr)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--corpus", choices=["gutenberg", "reuters", "all"],
                        default="all", help="which corpus to export (default: all)")
    parser.add_argument("--data-dir", default="data",
                        help="output directory (default: data)")
    args = parser.parse_args()

    data = Path(args.data_dir)
    data.mkdir(exist_ok=True)

    corpora = ["gutenberg", "reuters"] if args.corpus == "all" else [args.corpus]
    for name in corpora:
        out = data / f"{name}.jsonl"
        if out.exists():
            print(f"{out} already exists — delete it to re-export.", file=sys.stderr)
            continue
        if name == "gutenberg":
            fetch_gutenberg(out)
        elif name == "reuters":
            fetch_reuters(out)

    print("Done. Next step:", file=sys.stderr)
    for name in corpora:
        out = data / f"{name}.jsonl"
        if out.exists():
            print(f"  python3 scripts/tag-pos.py --file {out}", file=sys.stderr)


if __name__ == "__main__":
    main()
