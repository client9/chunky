# Tagging Pipeline

Given a string, the pipeline produces a slice of `Sentence`, each holding a slice of tagged `Token`.

Entry point: `tok.Parse(s string) []Sentence`

Each step after tokenization has the signature `func([]Token) []Token`, making it straightforward to test in isolation and to add new steps (emoji stripping, date merging, proper-noun-run merging) by inserting a function in the sequence.

---

## 1. Tokenization — `tok.Tokenize`

Splits the input string on Unicode whitespace. Each whitespace-delimited field becomes exactly one `Token{Word, Offset}`. No normalization, filtering, or tagging is applied. The byte offset of each token's first byte in the original string is recorded here and never changes.

---

## 2. Bracket stripping — `tok.StripBrackets`

Removes bracketed noise from the token stream. Three cases are handled:

1. **Whole-token bracket** (`[1]`, `[sic]`, `{x}`): token is removed entirely.
2. **Embedded numeric citation** (`word[8].`): the `[digits]` span is replaced with spaces of equal byte length, preserving the byte offsets of all surrounding characters. `SplitPunctuation` (step 4) then correctly locates the trailing punctuation via a last-non-space scan.
3. **Multi-token span** (`[critical section]`): tokens are buffered until the closing bracket and discarded. Unclosed spans pass through rather than consuming the rest of the stream.

The space-replacement approach for embedded citations is what allows a single function to handle both the embedded and standalone bracket cases without duplicating logic.

---

## 3. Text normalization — `tok.NormalizeText`

Applies `client9/typewriter` to each token's Word: curly/smart quotes → straight ASCII quotes, Unicode spaces → ASCII space. Runs before `SplitPunctuation` so that multi-character sequences are recognized on whole fields before punctuation is split off.

When `client9/demoji` stabilizes, emoji stripping will be added here as an additional per-token transform.

---

## 4. Punctuation splitting — `tok.SplitPunctuation`

Splits leading and trailing punctuation into separate tokens. For each token:

1. **Trim leading spaces** (produced by step 2 when an embedded citation started a field).
2. **Split leading `(`** into its own token.
3. **Find last non-space character** — scans right-to-left past any trailing spaces left by bracket replacement — and strip trailing `,.;:!?` into a separate token. The offset of that punctuation character is `Token.Offset + index`, which is exact because space-replacement preserves byte positions.
4. **Trim trailing spaces** from what remains.
5. **Split trailing `)`** into its own token.
6. **Respect `DottedAbbreviations`**: if the trimmed word ending in `.` is a known abbreviation (`Dr.`, `etc.`, `U.S.`), the dot stays attached.

---

## 5. Contraction splitting — `tok.SplitContractions`

Expands contraction tokens into (stem, suffix) pairs. No other punctuation knowledge here — this step consults the lexicon and contraction tables only.

- **Irregular forms** (`won't` → `will` + `n't`, `shan't` → `shall` + `n't`) via `ContractionNorm`.
- **Whole-word exceptions** — words in `AbbreviationTags` not in `ContractionNorm` stay whole (`ain't`, `o'clock`).
- **Auxiliary suffixes** (`'ll`, `'re`, `'ve`, `'m`, `'d`, `'s`) split at the apostrophe.
- **Negating `n't`**: if the stem-without-n is in the lexicon (`do`, `should`), the `n` moves to the suffix (`do` + `n't`); otherwise the `n` stays in the stem (`can` + `'t`, where `ca` is not a word).

---

## 6. Lexical compound merging — `tok.MergeLexical`

Scans left-to-right with longest-match for entries in `CompoundTags`. Matched sequences are replaced by a single token whose `Word` is the original surface form (space-joined: `"such as"`) and whose tag is the compound's UD tag (`"such as"` → ADP, `"in order to"` → PART). Offset is taken from the first token in the sequence.

Runs before `LexicalTag` so that individual words in a matched compound are not tagged unnecessarily.

---

## 7. Lexical tagging — `tok.LexicalTag`

For each untagged token (tokens already carrying candidates, e.g. compound tokens from step 6, are skipped), looks up the lowercase form:

1. **Compiled lexicon** (`wordtagmap`, ~35k entries generated from the Brown corpus + closed forms + hand-curated overrides). Assigns an ordered candidate tag slice.
2. **`AbbreviationTags` fallback** — runtime-editable map covering contraction suffixes (`'ll`, `n't`, …), titles, discourse abbreviations, and bibliographic forms not in the generated lexicon.

Tokens with no match in either map are left untagged (empty candidate slice).

---

## 8. Unknown-word tagging — `tok.TagUnknowns`

For each token still untagged, tries rules in order, stopping at the first match:

1. **Inflection** (`InflectionCandidates`): strips `-ing`, `-ed`, `-er`, `-est`, `-s`, `-es`, `-ies` and looks up the stem in the lexicon.
2. **Hyphen** (`HyphenCandidates`): splits at the last `-`; known adjectival suffixes (`-like`, `-free`, `-wide`) → ADJ unconditionally; otherwise looks up the final component in the lexicon or applies inflection/morph rules to it.
3. **Morphology** (`MorphCandidates`): prefix rules (`un-`, `re-`, `non-`, …) and suffix rules (`-tion`, `-ness`, `-ly`, `-ous`, …).
4. **Alpha fallback**: if the word is all Unicode letters and nothing else matched → NOUN (rule `unk:word`).

Numeric forms (integers, decimals, ordinals, decades) are tagged NUM by `MorphCandidates` before the above chain.

---

## 9. Sentence segmentation — `tok.Segment`

Splits the flat token stream into sentences and applies `RetagCapitalized` per sentence.

**Boundary detection** (`isBoundary`): a `.` `!` or `?` token is a sentence boundary unless:
- It is part of an ellipsis (adjacent `.` token).
- It follows a token in `AbbreviationTags` or `DottedAbbreviations`.
- It follows a single uppercase letter (middle initial).
- It sits between two PROPN tokens (PROPN `.` PROPN pattern for initials).

**`RetagCapitalized`** runs on each sentence independently so that `i == 0` correctly identifies the sentence-initial token:
- **`i == 0`** (sentence-initial): no change. Grammatical capitalization provides no information about the tag.
- **`i > 0`**, capitalized, known tag, not PRON: promote to PROPN. Handles proper nouns, Roman numerals, and initials that appear mid-sentence.

---

## Output

A `[]Sentence`, each containing:
- `Tokens []Token` — ordered token slice with `Word`, `Offset`, `Candidates []Tag`, `Rule`.
- `Offset int` — byte offset of the first token in the original string.

At this point each token carries an ordered candidate tag set. Disambiguation (context rules derived from corpus statistics) and chunking (MOD/HEAD assignment) are subsequent steps not yet implemented.
