# Tagging Pipeline

Given a string, the pipeline produces a slice of `Sentence`, each holding a slice of tagged `Token`.

---

## 1. Surface tokenization — `tok.surfaceTokenizeRaw`

Operates on the raw input string. Produces `rawToken{word, byteOffset}` pairs.

1. **Split on Unicode whitespace.** Each whitespace-delimited field becomes a candidate token. Byte offset into the original string is recorded at this point and never changes.
2. **Typographic normalization** via `client9/typewriter`: curly quotes → straight, em/en dashes → ASCII, Unicode spaces → space.
3. **Strip inline citations** (`[1]`, `[42]`): numeric bracket markers are removed from the right end of a field. Offsets of surrounding tokens are unaffected.
4. **Split leading `(`** into its own token.
5. **Strip trailing punctuation** (`,` `.` `:` `;` `!` `?`) into a separate token — *unless* the whole word (dot included) is a dotted abbreviation in `DottedAbbreviations` (e.g. `p.`, `Dr.`, `etc.`), in which case the dot stays attached.
6. **Split trailing `)`** into its own token.
7. **Split contractions** — `splitContractions` post-pass:
   - Irregular forms (`won't` → `will` + `n't`, `shan't` → `shall` + `n't`) via `ContractionNorm`.
   - Words in `AbbreviationTags` that are not in `ContractionNorm` stay whole (`ain't`).
   - Auxiliary suffixes (`'ll`, `'re`, `'ve`, `'m`, `'d`, `'s`) split at the apostrophe.
   - Negating `n't`: if the stem-without-n is in the lexicon (`do`, `should`), the `n` moves to the suffix (`do` + `n't`); otherwise the `n` stays in the stem (`can` + `'t`).

---

## 2. Lexical tagging — `tok.TagString`

For each raw token, look up the lowercase form:

1. **Compiled lexicon** (`wordtagmap`, ~35k entries generated from the Brown corpus + closed forms + hand-curated overrides). Assigns an ordered candidate tag slice.
2. **`AbbreviationTags` fallback** — runtime-editable map covering contraction suffixes (`'ll`, `n't`, …), titles, discourse abbreviations, and bibliographic forms not yet in the generated lexicon.

Tokens with no match in either map are left untagged (empty candidate slice).

---

## 3. Bracket filtering — `tok.FilterBrackets`

Removes bracketed noise spans from the token stream:

- Single-token forms (`[1]`, `[sic]`, `{x}`) are dropped.
- Multi-token spans (`[critical section]`) are buffered until the closing bracket and discarded.
- Unclosed brackets pass through rather than consuming the rest of the stream.

---

## 4. Compound merging — `tok.MergeCompounds`

Scans left-to-right with longest-match for entries in `CompoundTags`. Matched sequences are replaced by a single token whose word is the components joined by `_` and whose tag is the compound's UD tag (e.g. `such_as` → ADP, `in_order_to` → PART). Byte offset is taken from the first token in the sequence.

---

## 5. Unknown-word tagging — `tok.TagUnknowns`

For each token still untagged, tries rules in order, stopping at the first match:

1. **Inflection** (`InflectionCandidates`): strips `-ing`, `-ed`, `-er`, `-est`, `-s`, `-es`, `-ies` and looks up the stem in the lexicon.
2. **Hyphen** (`HyphenCandidates`): splits at the last `-`; known adjectival suffixes (`-like`, `-free`, `-wide`) → ADJ unconditionally; otherwise looks up the final component in the lexicon or applies inflection/morph rules to it.
3. **Morphology** (`MorphCandidates`): prefix rules (`un-`, `re-`, `non-`, …) and suffix rules (`-tion`, `-ness`, `-ly`, `-ous`, …).
4. **Alpha fallback**: if the word is all Unicode letters and nothing else matched → NOUN (rule `unk:word`).

Numeric forms (integers, decimals, ordinals, decades) are tagged NUM by `MorphCandidates` before the above chain.

---

## 6. Sentence segmentation — `tok.Segment`

Splits the flat token stream into sentences and applies `LexicalRetag` per sentence.

**Boundary detection** (`isBoundary`): a `.` `!` or `?` token is a sentence boundary unless:
- It is part of an ellipsis (adjacent `.` token).
- It follows a token in `AbbreviationTags` or `DottedAbbreviations`.
- It follows a single uppercase letter (middle initial).
- It sits between two PROPN tokens (PROPN `.` PROPN pattern for initials).

**`LexicalRetag`** runs on each sentence independently so that `i == 0` correctly identifies the sentence-initial token:
- **`i == 0`** (sentence-initial): no change. Grammatical capitalization provides no information about the tag.
- **`i > 0`**, capitalized, known tag, not PRON: promote to PROPN. Handles proper nouns, Roman numerals, and initials that appear mid-sentence.

---

## Output

A `[]Sentence`, each containing:
- `Tokens []Token` — ordered token slice with `Word`, `Offset`, `Canidates []Tag`, `Rule`.
- `Offset int` — byte offset of the first token in the original string.

At this point each token carries an ordered candidate tag set. Disambiguation (context rules derived from corpus statistics) and chunking (MOD/HEAD assignment) are subsequent steps not yet implemented.
