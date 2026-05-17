package tok

import "strings"

// DisambiguateOnce resolves the ADV/SCONJ ambiguity on "once".
//
// Frequency/temporal (ADV):  "happened once", "had once been", "once again"
// Subordinator (SCONJ):      "once the treaty was signed"
//
// Resolved cases:
//   - prev=AUX|ADV → ADV  ("had once been", "once again" when prev is adverb)
//   - next=ADV     → ADV  ("once more", "once again")
//   - next=ADP     → ADV  ("once upon a time", "once in a while")
//   - next=a/an    → ADV  ("once a month", "once an hour" — frequentive)
//   - next=VERB|AUX → SCONJ  ("once completed", "once approved")
//   - next=PRON    → SCONJ  ("once he arrived", "once they left")
//   - next=the     → SCONJ  ("once the treaty was signed")
//   - prev=VERB(pure) → ADV  ("it happened once.", "visited once more")
//   - next=PUNCT|CCONJ → ADV  ("once.", "once,", "once and for all")
func DisambiguateOnce(tokens []Token) []Token {
	for i := range tokens {
		disambiguateOnce(tokens, i)
	}
	return tokens
}

func disambiguateOnce(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADV) || !t.HasTag(TagSCONJ) {
		return
	}
	prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case prev.HasTag(TagAUX|TagADV) || next.HasTag(TagADV):
		resolve = TagADV
	case next.HasTag(TagSCONJ | TagADP):
		resolve = TagADV // "once upon a time", "once in a while", "once before"
	case strings.ToLower(next.Word) == "a" || strings.ToLower(next.Word) == "an":
		resolve = TagADV // "once a month", "once an hour" (frequentive)
	case next.HasTag(TagVERB | TagAUX):
		resolve = TagSCONJ // "once completed", "once approved"
	case next.HasTag(TagPRON):
		resolve = TagSCONJ // "once he arrived", "once they left"
	case next.Word == "the":
		resolve = TagSCONJ // "once the treaty was signed"
	case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
		resolve = TagADV // "it happened once.", "visited once more"
	case next.HasTag(TagPUNCT | TagCCONJ):
		resolve = TagADV // "once.", "once,", "once and for all"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+so"
	}
}
