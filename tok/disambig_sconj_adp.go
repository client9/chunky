package tok

import "strings"

// DisambiguateSCONJasADP resolves words normally tagged SCONJ that function
// as prepositions when followed by a nominal complement.
//
//   - since:   ADP when next is DET|NOUN|PROPN|PRON|NUM|ADJ ("since the crash",
//     "since then" stays ambiguous)
//   - despite: always ADP — it has no genuine SCONJ use in prose
//   - upon:    always ADP — Brown corpus over-tagged it as SCONJ
func DisambiguateSCONJasADP(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagSCONJ) {
			continue
		}
		lw := strings.ToLower(t.Word)
		var resolve Tag
		switch lw {
		case "since":
			next := tokenAt(tokens, i+1)
			if next.HasTag(TagDET | TagNOUN | TagPROPN | TagPRON | TagNUM | TagADJ) {
				resolve = TagADP
			}
		case "despite", "upon":
			resolve = TagADP
		default:
			continue
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+sconj-adp"
		}
	}
	return tokens
}
