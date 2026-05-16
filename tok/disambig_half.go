package tok

import "strings"

// DisambiguateHalf resolves the ADJ/DET/NOUN ambiguity on "final", "half", and "individual".
//
// "final" lexicon: ADJ|DET|NOUN — almost always ADJ (87%)
//   - next=NOUN|ADJ|PROPN → ADJ  ("the final answer", "the final chapter")
//
// "half" lexicon: ADJ|DET|NOUN
//   - next=DET|PRON → DET   ("half the city", "half a mile")
//   - next=ADP      → NOUN  ("the first half of", "half of the season")
//
// "individual" lexicon: ADJ|DET|NOUN
//   - next=NOUN|ADJ|PROPN → ADJ   ("individual rights", "individual cases")
//   - next=PART            → NOUN  ("each individual to decide")
func DisambiguateHalf(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagADJ) || !t.HasTag(TagNOUN) {
			continue
		}
		lw := strings.ToLower(t.Word)
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch lw {
		case "final":
			switch {
			case next.HasTag(TagNOUN | TagADJ | TagPROPN | TagNUM):
				resolve = TagADJ
			case next.HasTag(TagPUNCT):
				resolve = TagADJ // "the final." — overwhelmingly ADJ (82%)
			}
		case "half":
			prev := tokenAt(tokens, i-1)
			switch {
			case next.HasTag(TagDET | TagPRON):
				resolve = TagDET
			case next.HasTag(TagADP | TagAUX):
				resolve = TagNOUN // "half of", "the second half was"
			case prev.HasTag(TagNUM | TagADJ):
				resolve = TagNOUN // "the second half", "the first half"
			case next.HasTag(TagPUNCT | TagCCONJ):
				resolve = TagNOUN // "in the second half.", "halftime"
			}
		case "individual":
			switch {
			case next.HasTag(TagNOUN | TagADJ | TagPROPN):
				resolve = TagADJ
			case next.HasTag(TagPART | TagADP | TagAUX):
				resolve = TagNOUN // "each individual to", "each individual of", "each individual was"
			}
		default:
			continue
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+half"
		}
	}
	return tokens
}
