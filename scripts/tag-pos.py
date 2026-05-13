#!/usr/bin/env python3
"""
Read a JSONL file produced by fetch-wiki-dump.py and emit a plain-text
POS-tagged corpus: one sentence per line, each token tagged word/TAG.

Example:
  The/DET dog/NOUN runs/VERB ./PUNCT

Output filename is derived from the source filename, e.g.:
  enwiki-20260401-multistream1.pos.txt

Usage:
  python3 scripts/tag-pos.py --file data/enwiki-20260401-multistream1.jsonl
  python3 scripts/tag-pos.py --file data/enwiki-20260401-multistream1.jsonl --limit 5000
"""

import argparse
import json
import re
import sys
from pathlib import Path

import spacy


def output_path(source_path, data_dir=None):
    p = Path(source_path)
    stem = p.stem  # strips one suffix (.jsonl → base name)
    out_dir = Path(data_dir) if data_dir else p.parent
    return out_dir / f"{stem}.pos.txt"


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", required=True, help="input .jsonl file")
    parser.add_argument("--limit", type=int, default=0,
                        help="max articles to process (0 = all)")
    parser.add_argument("--model", default="en_core_web_sm",
                        help="spaCy model to use (default: en_core_web_sm)")
    args = parser.parse_args()

    src = Path(args.file)
    if not src.exists():
        print(f"Error: file not found: {src}", file=sys.stderr)
        sys.exit(1)

    out = output_path(src)
    if out.exists():
        print(f"Output {out} already exists — delete it first to re-tag.", file=sys.stderr)
        sys.exit(1)

    print(f"Loading spaCy model '{args.model}' ...", file=sys.stderr)
    nlp = spacy.load(args.model)

    print(f"Tagging → {out}", file=sys.stderr)

    count = 0
    sent_count = 0
    with open(src) as fin, open(out, "w") as fout:
        for line in fin:
            line = line.strip()
            if not line:
                continue
            record = json.loads(line)
            text = record.get("extract", "")
            if not text:
                continue

            doc = nlp(text)
            for sent in doc.sents:
                tokens = [
                    f"{tok.text}/{tok.pos_}"
                    for tok in sent
                    if not tok.is_space
                ]
                if tokens:
                    fout.write(" ".join(tokens) + "\n")
                    sent_count += 1

            count += 1
            if count % 1000 == 0:
                print(f"  {count} articles, {sent_count} sentences", file=sys.stderr)
            if args.limit and count >= args.limit:
                break

    print(f"Done. {count} articles, {sent_count} sentences → {out}", file=sys.stderr)


if __name__ == "__main__":
    main()
