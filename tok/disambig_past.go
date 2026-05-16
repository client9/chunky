package tok

import "strings"

// DisambiguatePast resolves the ADJ/ADP/NOUN ambiguity on "past" and "pro".
//
// "past" (ADJ 45%, ADP 30%, NOUN 25%):
//   - next=NOUN|NUM|ADJ → ADJ  ("past performance", "past 10 years", "past president")
//   - next=PUNCT|AUX    → NOUN ("in the past.", "the past was")
//   - next=DET|PROPN    → ADP  ("past the building", "past Paris")
//
// "pro" (overwhelmingly ADJ in corpus):
//   - next=ADJ|NOUN|PROPN → ADJ  ("pro rata", "pro wrestler", "pro sports")
func DisambiguatePast(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagADJ) {
			continue
		}
		lw := strings.ToLower(t.Word)
		switch lw {
		case "past", "pro":
		default:
			continue
		}
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch lw {
		case "past":
			if !t.HasTag(TagADP) && !t.HasTag(TagNOUN) {
				continue
			}
			switch {
			case next.HasTag(TagNOUN | TagNUM | TagADJ):
				resolve = TagADJ
			case next.HasTag(TagPUNCT | TagAUX):
				resolve = TagNOUN
			case next.HasTag(TagDET | TagPROPN):
				resolve = TagADP
			}
		case "pro":
			if next.HasTag(TagADJ | TagNOUN | TagPROPN) {
				resolve = TagADJ
			}
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+past"
		}
	}
	return tokens
}
