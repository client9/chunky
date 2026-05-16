package tok

import "strings"

// DisambiguatePrepositions resolves the ADP/ADV ambiguity on spatial prepositions
// that double as adverbs: about, around, below, behind, out.
//
// ADP when introducing a complement NP:
//   - next=DET|NOUN|PRON|PROPN|NUM → ADP  (93–100% in corpus)
//
// ADV when after a verb with no complement:
//   - around/below + next=PUNCT  → ADV  ("turned around.", "shown below.")
//   - behind + prev=VERB(pure)   → ADV  ("fell behind")
//
// "about" is ADP in almost all non-numeric contexts; numeric "about 3" is left
// ambiguous (76.7% ADV, 23.3% ADP — genuine approximator vs preposition).
func DisambiguatePrepositions(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagADP) || !t.HasTag(TagADV) {
			continue
		}
		lw := strings.ToLower(t.Word)
		switch lw {
		case "about", "around", "below", "behind", "out":
		default:
			continue
		}
		prev := tokenAt(tokens, i-1)
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch lw {
		case "about":
			switch {
			case next.HasTag(TagDET | TagNOUN | TagPRON | TagPROPN | TagSCONJ | TagADJ | TagVERB | TagPART):
				resolve = TagADP // "about the issue", "about to leave"
			case next.HasTag(TagPUNCT):
				resolve = TagADV // "that's what it's about."
			}
		case "around":
			switch {
			case next.HasTag(TagDET | TagNOUN | TagPRON | TagPROPN | TagADJ | TagNUM):
				resolve = TagADP
			case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
				resolve = TagADV // "looked around", "turned around"
			case next.HasTag(TagPUNCT | TagCCONJ | TagADV):
				resolve = TagADV // "turned around.", "around and about"
			}
		case "below":
			switch {
			case next.HasTag(TagDET | TagNOUN | TagPRON | TagPROPN | TagNUM):
				resolve = TagADP
			case next.HasTag(TagPUNCT | TagCCONJ):
				resolve = TagADV // "shown below.", "below and above"
			}
		case "behind":
			switch {
			case next.HasTag(TagDET | TagNOUN | TagPROPN | TagPRON):
				resolve = TagADP
			case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
				resolve = TagADV // "fell behind", "left behind"
			case next.HasTag(TagPUNCT | TagCCONJ | TagADV):
				resolve = TagADV // "fell behind.", "behind and before"
			}
		case "out":
			switch {
			case next.HasTag(TagDET | TagNOUN | TagADJ | TagPROPN | TagPRON):
				resolve = TagADP // "out the door", "out of bounds"
			case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
				resolve = TagADV // "burned out", "went out", "carried out"
			case next.HasTag(TagPUNCT | TagCCONJ | TagADV):
				resolve = TagADV // "went out.", "out and about", "out there"
			}
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+prep"
		}
	}
	return tokens
}
