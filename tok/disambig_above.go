package tok

import "strings"

// DisambiguateAbove resolves {ADJ,ADP,ADV} location words: outside, above, inside.
//
// "outside" (prenominal ADJ common):
//   - next=NOUN|ADJ|PROPN → ADJ  ("outside world", "outside influence")
//   - next=DET|PRON|NUM   → ADP  ("outside the building", "outside us")
//   - prev=VERB (pure)    → ADV  ("went outside")
//
// "above", "inside":
//   - next=DET|NOUN|PROPN|NUM|PRON|ADJ → ADP  ("above sea level", "above average")
//   - next=PUNCT          → ADV  ("mentioned above.")
//   - prev=VERB (pure)    → ADV  ("rose above")
func DisambiguateAbove(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagADP) || !t.HasTag(TagADV) {
			continue
		}
		lw := strings.ToLower(t.Word)
		switch lw {
		case "outside", "above", "inside":
		default:
			continue
		}
		prev := tokenAt(tokens, i-1)
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch lw {
		case "outside":
			switch {
			case next.HasTag(TagNOUN | TagADJ | TagPROPN):
				resolve = TagADJ // "outside world", "outside influence"
			case next.HasTag(TagDET | TagPRON | TagNUM):
				resolve = TagADP // "outside the building", "outside us"
			case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
				resolve = TagADV // "went outside"
			case next.HasTag(TagPUNCT):
				resolve = TagADV // "stayed outside."
			}
		default: // above, inside
			switch {
			case next.HasTag(TagDET | TagNOUN | TagPROPN | TagNUM | TagPRON | TagADJ):
				resolve = TagADP
			case next.HasTag(TagPUNCT):
				resolve = TagADV // "mentioned above.", "stayed inside."
			case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
				resolve = TagADV
			}
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+above"
		}
	}
	return tokens
}
