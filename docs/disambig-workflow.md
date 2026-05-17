# Word-Specific Disambiguation Workflow

Word-specific disambiguators live in `tok/disambig_*.go` and handle two cases that the generated context rules (`rules_gen.go`) cannot:

1. **Triple-tagged words** — `rules_gen.go` only fires on 2-tag tokens. A word like `up {ADP,ADV,NOUN}` needs a dedicated function.
2. **Lexically idiosyncratic behavior** — words where tag identity matters, not just the surrounding tag pattern (e.g. `as`, `will`, `there`).

---

## Step 1 — Find words worth targeting

```bash
# By tag combination — best starting point; shows which ambiguity class has the most instances
go run ./cmd/ambig -tags data/*.pos.jsonl | head -20

# By word — drill into a specific class after the above
go run ./cmd/ambig data/*.pos.jsonl | head -60
```

Prioritize triple-tagged words first (no context rules can touch them), then high-count 2-tag pairs that the generated rules aren't resolving.

---

## Step 2 — Analyze context patterns

Use the helper script (see below) or run manually:

```bash
word=up
go run ./cmd/lexrules -feat word+nexttag -sort-count data/*.pos.jsonl | grep -E "^\s+[0-9]+\. ${word}\+"
go run ./cmd/lexrules -feat word+prevtag -sort-count data/*.pos.jsonl | grep -E "^\s+[0-9]+\. ${word}\+"
```

**Reading the output:**  Each row is `word+CONTEXT (n=COUNT): TAG: PCT%, ...`

Write a rule only when:
- n ≥ 200 (enough evidence)
- Dominant tag ≥ 85% (confident)
- The pattern has a linguistic explanation (not just a corpus artifact)

Watch for interactions: `word+nexttag` and `word+prevtag` are single-slot features. If the two features suggest conflicting resolutions for the same token (e.g. `prev=VERB→ADP` and `next=ADP→ADV` both apply), use the sandwich feature or leave it unresolved rather than guess.

```bash
# Sandwich (prev+next together) — use when single-slot features conflict
go run ./cmd/lexrules -feat prevtag+nexttag -sort-count data/*.pos.jsonl | grep -E "^\s+[0-9]+\. [A-Z]+\+[A-Z]+" | head -30
```

---

## Step 3 — Write the rule

Model new files on an existing one, e.g. `tok/disambig_up.go`. The pattern:

```go
func DisambiguateX(tokens []Token) []Token {
    for i, t := range tokens {
        if strings.ToLower(t.Word) != "x" {
            continue
        }
        if !t.HasTag(TagFOO) {   // guard: skip if already resolved away
            continue
        }
        prev := tokenAt(tokens, i-1)
        next := tokenAt(tokens, i+1)
        var resolve Tag
        switch {
        case next.HasTag(TagDET | TagNOUN):
            resolve = TagADP   // "x the ...", "x a ..."
        // ...
        }
        if resolve != 0 {
            tokens[i].Tags = resolve
            tokens[i].Rule = t.Rule + "+x"
        }
    }
    return tokens
}
```

Key helpers:
- `resolvedAs(tok, tag)` — tok is fully resolved AND has that tag (use for prev/next that must be unambiguous)
- `tok.HasTag(tag)` — tok has that tag among its candidates (allows ambiguous neighbors)
- `tokenAt(tokens, i)` — safe; returns a zero Token at boundaries

Use `resolvedAs` for the neighbor when a false positive would be worse than leaving it unresolved. Use `HasTag` when you want to fire even if the neighbor is still ambiguous.

Register the new function in `tok/disambig_words.go`. Order matters only if one pass produces input for another; otherwise append near related words.

---

## Step 4 — Verify and test

```bash
# Check a specific sentence
echo "they picked up the package ." | go run ./cmd/tok

# Run corpus ambiguity count before/after
go run ./cmd/ambig data/*.pos.jsonl | grep -E "^\s+[0-9]+\. up\b"
```

Write a `tok/disambig_X_test.go` with cases covering each branch. Use sentences where you've confirmed the expected tag with the CLI above. Watch for compounds: `MergeLexical` runs before all disambiguators, so `as well`, `such as`, `as long as`, etc. never reach word-specific rules.

---

## Helper script

`scripts/disambig-analyze.sh <word>` — runs both feature queries and filters to the target word:

```bash
#!/usr/bin/env bash
set -euo pipefail
word="${1:?usage: disambig-analyze.sh <word>}"
corpus="data/*.pos.jsonl"

echo "=== next-tag patterns for '$word' ==="
go run ./cmd/lexrules -feat word+nexttag -sort-count $corpus | grep -E "^\s+[0-9]+\. ${word}\+"

echo ""
echo "=== prev-tag patterns for '$word' ==="
go run ./cmd/lexrules -feat word+prevtag -sort-count $corpus | grep -E "^\s+[0-9]+\. ${word}\+"
```

Save as `scripts/disambig-analyze.sh`, `chmod +x`, then:

```bash
bash scripts/disambig-analyze.sh up
bash scripts/disambig-analyze.sh long
```
