package tok

import "strings"

// DisambiguatePrepositions resolves the ADP/ADV ambiguity on spatial prepositions
// that double as adverbs: about, around, below, behind, out.
//
// ADP when introducing a complement NP:
//   - next=DET|NOUN|PRON|PROPN|NUM → ADP  (93–100% in corpus)
//
// ADP when following a resolved verb (phrasal particle):
//   - prev=resolved-VERB + next≠PUNCT → ADP  ("look out X", "fall behind X")
//
// ADV when after a verb with no complement:
//   - around/behind/out + prev=VERB(pure) → ADV  ("looked around", "fell behind")
//   - around/below + next=PUNCT → ADV  ("turned around.", "shown below.")
//
// ADV approximator:
//   - about/around + next=NUM → ADV  ("about 50", "around 100 years")
func DisambiguatePrepositions(tokens []Token) []Token {
	for i := range tokens {
		disambiguatePrepositions(tokens, i)
	}
	return tokens
}

func disambiguatePrepositions(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADP) || !t.HasTag(TagADV) {
		return
	}
	lw := strings.ToLower(t.Word)
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch lw {
	case "about":
		switch {
		case resolvedAs(prev, TagVERB) && !next.HasTag(TagPUNCT):
			resolve = TagADP // "talk about X", "ask about Y" — verb-particle
		case next.HasTag(TagDET | TagNOUN | TagPRON | TagPROPN | TagSCONJ | TagADJ | TagVERB | TagPART):
			resolve = TagADP // "about the issue", "about to leave"
		case next.HasTag(TagPUNCT):
			resolve = TagADV // "that's what it's about."
		case next.HasTag(TagNUM):
			resolve = TagADV // "about 50" — approximator
		}
	case "around":
		switch {
		case next.HasTag(TagNUM):
			resolve = TagADV // "around 100 years" — approximator
		case resolvedAs(prev, TagVERB) && !next.HasTag(TagPUNCT):
			resolve = TagADP // "move around X" — verb-particle
		case next.HasTag(TagDET | TagNOUN | TagPRON | TagPROPN | TagADJ):
			resolve = TagADP
		case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
			resolve = TagADV // "looked around", "turned around"
		case next.HasTag(TagPUNCT | TagCCONJ | TagADV):
			resolve = TagADV // "turned around.", "around and about"
		}
	case "below":
		switch {
		case resolvedAs(prev, TagVERB) && !next.HasTag(TagPUNCT):
			resolve = TagADP // "look below X" — verb-particle
		case next.HasTag(TagDET | TagNOUN | TagPRON | TagPROPN | TagNUM):
			resolve = TagADP
		case next.HasTag(TagPUNCT | TagCCONJ):
			resolve = TagADV // "shown below.", "below and above"
		}
	case "behind":
		switch {
		case resolvedAs(prev, TagVERB) && !next.HasTag(TagPUNCT):
			resolve = TagADP // "look behind X" — verb-particle
		case next.HasTag(TagDET | TagNOUN | TagPROPN | TagPRON):
			resolve = TagADP
		case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
			resolve = TagADV // "fell behind", "left behind"
		case next.HasTag(TagPUNCT | TagCCONJ | TagADV):
			resolve = TagADV // "fell behind.", "behind and before"
		}
	case "out":
		switch {
		case resolvedAs(prev, TagVERB) && !next.HasTag(TagPUNCT):
			resolve = TagADP // "burned out X" — verb-particle
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
