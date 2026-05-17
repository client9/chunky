package tok

import "strings"

// DisambiguateRight resolves "right" ({ADJ,ADV,NOUN}).
// No single tag dominates across all contexts; next=PUNCT (45/27/26 split) is left unresolved.
//
//   - next=NOUN|PROPN       → ADJ   ("right time", "right place")        95%
//   - next=PART             → NOUN  ("right to vote", "right to remain")  96%
//   - next=ADV              → ADV   ("right now", "right away")           80%
//   - prev=AUX              → ADJ   ("is right", "are right")             87%
func DisambiguateRight(tokens []Token) []Token {
	for i, t := range tokens {
		if strings.ToLower(t.Word) != "right" {
			continue
		}
		if !t.HasTag(TagADJ) {
			continue
		}
		prev := tokenAt(tokens, i-1)
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch {
		case next.HasTag(TagNOUN | TagPROPN):
			resolve = TagADJ
		case next.HasTag(TagPART):
			resolve = TagNOUN // "right to vote"
		case resolvedAs(next, TagADV):
			resolve = TagADV // "right now", "right away", "right here"
		case resolvedAs(prev, TagAUX):
			resolve = TagADJ // "is right", "are right", "was right"
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+right"
		}
	}
	return tokens
}
