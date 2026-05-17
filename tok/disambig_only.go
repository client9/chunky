package tok

import "strings"

// DisambiguateOnly resolves the ADJ/ADV/DET ambiguity on "only" and "little".
//
// only:
//   - prev=DET + next=NOUN|PROPN|PRON|NUM → ADJ  ("the only person", "the only one")
//   - next=VERB|AUX|ADV|ADP|SCONJ|PART|PRON|DET|NUM|PUNCT → ADV
//   - next=ADJ (unambiguous) → ADV
//
// little:
//   - next=ADJ|ADV → ADV  ("a little tired")
//   - next=NOUN|PROPN → DET  ("little time")
func DisambiguateOnly(tokens []Token) []Token {
	for i := range tokens {
		disambiguateOnly(tokens, i)
	}
	return tokens
}

func disambiguateOnly(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
		return
	}
	lw := strings.ToLower(t.Word)
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch lw {
	case "only":
		switch {
		case resolvedAs(prev, TagDET) && next.HasTag(TagNOUN|TagPROPN|TagPRON|TagNUM):
			resolve = TagADJ
		case next.HasTag(TagVERB | TagAUX | TagADV | TagADP | TagSCONJ | TagPART | TagPRON | TagDET | TagNUM | TagPUNCT):
			resolve = TagADV
		case next.HasTag(TagADJ) && !next.HasTag(TagNOUN):
			resolve = TagADV
		}
	case "little":
		switch {
		case next.HasTag(TagADJ | TagADV):
			resolve = TagADV
		case next.HasTag(TagNOUN | TagPROPN):
			resolve = TagDET
		}
	default:
		return
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adv-noun"
	}
}
