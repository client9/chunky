package tok

import "strings"

// DisambiguateWellBack resolves the ADJ|ADV|NOUN|VERB ambiguity on "back" and "well".
//
// back:
//   - prev=VERB            → ADV  ("went back", "came back", "step back")
//   - next=NOUN, prev≠VERB → ADJ  ("back door", "back seat")
//
// well:
//   - prev=VERB            → ADV  ("did well", "worked well")
//   - prev=DET             → NOUN ("the well", "an oil well")
func DisambiguateWellBack(tokens []Token) []Token {
	for i, t := range tokens {
		lw := strings.ToLower(t.Word)
		switch lw {
		case "back":
			if !t.HasTag(TagADV) || !t.HasTag(TagNOUN) {
				continue
			}
			prev := tokenAt(tokens, i-1)
			next := tokenAt(tokens, i+1)
			switch {
			case prev.HasTag(TagVERB):
				tokens[i].Tags = TagADV
				tokens[i].Rule = t.Rule + "+back"
			case next.HasTag(TagNOUN|TagADJ|TagPROPN) && !prev.HasTag(TagVERB):
				tokens[i].Tags = TagADJ
				tokens[i].Rule = t.Rule + "+back"
			}
		case "well":
			if !t.HasTag(TagADV) || !t.HasTag(TagNOUN) {
				continue
			}
			prev := tokenAt(tokens, i-1)
			switch {
			case prev.HasTag(TagVERB):
				tokens[i].Tags = TagADV
				tokens[i].Rule = t.Rule + "+well"
			case resolvedAs(prev, TagDET):
				tokens[i].Tags = TagNOUN
				tokens[i].Rule = t.Rule + "+well"
			}
		}
	}
	return tokens
}
