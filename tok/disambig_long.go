package tok

import "strings"

// DisambiguateLong resolves "long" ({ADJ,ADV}).
//
//   - next=NOUN|ADJ|PROPN|NUM (unambiguous) → ADJ  ("a long time", "long road")
//   - prev=DET                               → ADJ  ("the long", "a long")
//   - prev=VERB|AUX                          → ADV  ("waited long", "how long has")
//   - prev=PRON|NOUN|PROPN + next=PUNCT      → ADV  ("she knew best.", "he waited long.")
//   - next=resolved-VERB|AUX                 → ADV  ("how long has", "best served")
//   - next=ADP|SCONJ                         → ADV  ("long before", "as long as")
//   - next=AUX|ADV                           → ADV  ("how long will", "so long ago")
//   - next=VERB (not AUX)                    → ADV  ("how long waited")
func DisambiguateLong(tokens []Token) []Token {
	for i := range tokens {
		disambiguateLong(tokens, i)
	}
	return tokens
}

func disambiguateLong(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagNOUN|TagADJ|TagPROPN|TagNUM) && !next.HasTag(TagVERB|TagAUX):
		resolve = TagADJ
	case strings.ToLower(prev.Word) == "no":
		resolve = TagADV // "no longer"
	case resolvedAs(prev, TagDET):
		resolve = TagADJ
	case prev.HasTag(TagVERB | TagAUX):
		resolve = TagADV
	case prev.HasTag(TagPRON|TagNOUN|TagPROPN) && next.HasTag(TagPUNCT):
		resolve = TagADV // "he waited long.", "she knew best."
	case resolvedAs(next, TagVERB) || resolvedAs(next, TagAUX):
		resolve = TagADV // "how long has", "best served cold"
	case next.HasTag(TagADP | TagSCONJ):
		resolve = TagADV // "long before", "as long as", "best of all"
	case next.HasTag(TagAUX | TagADV):
		resolve = TagADV // "how long will", "so long ago"
	case next.HasTag(TagVERB) && !next.HasTag(TagAUX):
		resolve = TagADV // "how long waited", "not long after"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+long"
	}
}
