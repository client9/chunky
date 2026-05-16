package tok

import "strings"

// DisambiguateSpatial resolves the ADP/ADV ambiguity on common spatial prepositions.
//
// All words: prev=VERB → ADP (verb-particle / phrasal-preposition: "look out", "fall behind")
// "about", "around": next=NUM → ADV (approximation: "about 50", "around 100 years")
func DisambiguateSpatial(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagADP) || !t.HasTag(TagADV) {
			continue
		}
		lw := strings.ToLower(t.Word)
		prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
		var resolve Tag
		switch lw {
		case "out", "below", "behind":
			if resolvedAs(prev, TagVERB) && !next.HasTag(TagPUNCT) {
				resolve = TagADP
			}
		case "about", "around":
			switch {
			case resolvedAs(prev, TagVERB) && !next.HasTag(TagPUNCT):
				resolve = TagADP
			case next.HasTag(TagNUM):
				resolve = TagADV
			}
		default:
			continue
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+spatial"
		}
	}
	return tokens
}
