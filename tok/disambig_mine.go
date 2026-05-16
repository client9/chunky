package tok

import "strings"

// DisambiguateMine resolves the NOUN/PRON ambiguity on "mine", "U"/"u".
//
// mine — PRON (possessive) or NOUN (excavation):
//   - prev=AUX|VERB → PRON  ("is mine", "was mine")
//   - prev=DET       → NOUN  ("the mine", "a mine")
//   - next=NOUN|PROPN → NOUN  ("mine shaft")
//   - next=ADP       → NOUN  ("mine of gold")
//
// U/u — almost always NOUN (letter, shape, abbreviation):
//   - prev=DET|ADJ|NUM → NOUN  ("a U shape", "3 U bolts")
//   - next=NOUN|PROPN   → NOUN  ("U shape", "U bolt")
//   - next=PUNCT         → NOUN  (letter reference: "vitamin U.")
func DisambiguateMine(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagNOUN) || !t.HasTag(TagPRON) {
			continue
		}
		lw := strings.ToLower(t.Word)
		prev := tokenAt(tokens, i-1)
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch lw {
		case "mine":
			switch {
			case prev.HasTag(TagAUX | TagVERB):
				resolve = TagPRON
			case resolvedAs(prev, TagDET):
				resolve = TagNOUN
			case next.HasTag(TagNOUN | TagPROPN):
				resolve = TagNOUN
			case next.HasTag(TagADP):
				resolve = TagNOUN
			}
		case "u":
			switch {
			case prev.HasTag(TagDET | TagADJ | TagNUM):
				resolve = TagNOUN
			case next.HasTag(TagNOUN | TagPROPN):
				resolve = TagNOUN // "U shape", "U bolt"
			case next.HasTag(TagPUNCT):
				resolve = TagNOUN // "vitamin U."
			}
		default:
			continue
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+mine"
		}
	}
	return tokens
}
