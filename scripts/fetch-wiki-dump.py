#!/usr/bin/env python3
"""
Parse a Wikipedia multistream XML dump, extracting article text
as JSONL for use with analyze-patterns.py.

Output filename is derived from the source filename, e.g.:
  enwiki-20260401-multistream1.jsonl

Usage:
  python3 scripts/fetch-wiki-dump.py --file enwiki-...-multistream1.xml.bz2
  python3 scripts/fetch-wiki-dump.py --file enwiki-...-multistream1.xml.bz2 --limit 5000
"""

import argparse
import bz2
import html
import json
import re
import sys
from pathlib import Path

def output_path(source_name, data_dir="data"):
    """Derive a clean output filename from the dump filename."""
    # enwiki-20260401-pages-articles-multistream1.xml-p1p41242.bz2
    # → enwiki-20260401-multistream1.jsonl
    m = re.search(r'(enwiki-\d+).*?(multistream\d+)', source_name)
    if m:
        stem = f"{m.group(1)}-{m.group(2)}"
    else:
        stem = Path(source_name).stem.split('.')[0]
    return Path(data_dir) / f"{stem}.jsonl"

# Strip wikitext markup, leaving plain prose.
def clean(text):
    text = text.replace('&nbsp;', ' ')
    text = html.unescape(text)               # &lt;ref&gt; → <ref>, etc.
    text = text.replace('&nbsp;', ' ')       # catch &amp;nbsp; → &nbsp; after unescape
    text = re.sub(r'<math[^>]*>.*?</math>', '', text, flags=re.DOTALL)  # LaTeX math blocks
    text = re.sub(r'<ref[^>]*/>', '', text)  # self-closing <ref />
    text = re.sub(r'<ref[^>]*>.*?</ref>', '', text, flags=re.DOTALL)  # <ref>citation</ref>
    # Strip {{ }} templates iteratively to handle nesting (infoboxes, etc.)
    while True:
        new = re.sub(r'\{\{[^{}]*\}\}', '', text)
        if new == text:
            break
        text = new
    text = re.sub(r'^\[\[(?:File|Image):[^\n]*$', '', text, flags=re.IGNORECASE | re.MULTILINE)
    text = re.sub(r'\[\[(?:[^|\]]*\|)?([^\]]+)\]\]', r'\1', text)
    text = re.sub(r'\[https?://[^\s\]]+[^\]]*\]', '', text)
    text = re.sub(r"'{2,}", '', text)
    text = re.sub(r'<[^>]*>', '', text)      # strip all HTML/XML tags after unescaping
    text = re.sub(r'==+[^=]+=+', '', text)
    text = re.sub(r'[|!][^\n]*', '', text)  # strip table cells (| data, ! header)
    text = re.sub(r'\(\s*[,;:\s]+', '(', text)  # strip leading punctuation artifacts inside parens
    text = re.sub(r'[,;:\s]+\)', ')', text)      # strip trailing punctuation artifacts inside parens
    text = re.sub(r'\([,;\s]*\)', '', text)      # remove empty parens left by stripped templates
    text = re.sub(r'\n{3,}', '\n\n', text)  # collapse excessive blank lines
    text = text.replace('–', '-').replace('—', '--')
    text = re.sub(r'[ \t]+', ' ', text)
    return text.strip()

def parse_dump(stream):
    """Yield (title, text) pairs from an XML dump stream."""
    title = None
    in_text = False
    buf = []

    for raw in stream:
        line = raw.decode('utf-8', errors='replace')

        m = re.search(r'<title>(.*?)</title>', line)
        if m:
            title = m.group(1)
            continue

        if '<text' in line:
            in_text = True
            m = re.search(r'<text[^>]*>(.*)', line, re.DOTALL)
            if m:
                buf.append(m.group(1))
            continue

        if in_text:
            if '</text>' in line:
                buf.append(line[:line.index('</text>')])
                in_text = False
                text = ''.join(buf).strip()
                buf = []
                if title and text and not text.startswith('#REDIRECT'):
                    summary = text.split('==')[0]
                    yield title, clean(summary)
                title = None
            else:
                buf.append(line)

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--file', required=True,
                        help='local .bz2 dump file')
    parser.add_argument('--limit', type=int, default=0,
                        help='max articles to extract (0 = all)')
    args = parser.parse_args()

    path = Path(args.file)
    if not path.exists():
        print(f"Error: file not found: {path}", file=sys.stderr)
        sys.exit(1)
    raw, source_name = open(path, 'rb'), path.name

    out = output_path(source_name)
    out.parent.mkdir(parents=True, exist_ok=True)
    if out.exists():
        print(f"Output {out} already exists — delete it first to re-extract.", file=sys.stderr)
        sys.exit(1)

    print(f"Extracting → {out}", file=sys.stderr)

    count = 0
    with raw:
        stream = bz2.open(raw, 'rb')
        with open(out, 'w') as f:
            for title, text in parse_dump(stream):
                if not text:
                    continue
                f.write(json.dumps({"title": title, "extract": text}) + "\n")
                count += 1
                if count % 1000 == 0:
                    print(f"  {count} articles", file=sys.stderr)
                if args.limit and count >= args.limit:
                    break

    print(f"Done. {count} articles → {out}", file=sys.stderr)

if __name__ == '__main__':
    main()
