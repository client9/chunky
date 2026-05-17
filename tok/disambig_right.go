package tok

// DisambiguateRight resolves "right" ({ADJ,ADV,NOUN}).
// No single tag dominates across all contexts.
//
//   - next=NOUN|ADJ|PROPN   → ADJ   ("right time", "right place", "right decision")
//   - next=PART              → NOUN  ("right to vote", "right to remain")
//   - prev=resolved-AUX      → ADJ   ("is right", "are right")
//   - next=ADV|DET|PRON      → ADV   ("right now", "right away", "right here")
//   - prev=VERB(not AUX)     → ADV   ("turned right", "veered right")
//   - prev=DET + no noun     → NOUN  ("the right (to vote)")
func DisambiguateRight(tokens []Token) []Token {
	for i := range tokens {
		disambiguateRight(tokens, i)
	}
	return tokens
}

func disambiguateRight(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagNOUN | TagADJ | TagPROPN):
		resolve = TagADJ
	case next.HasTag(TagPART):
		resolve = TagNOUN // "right to vote"
	case resolvedAs(prev, TagAUX):
		resolve = TagADJ // "is right", "are right", "was right"
	case next.HasTag(TagADV | TagDET | TagPRON):
		resolve = TagADV // "right now", "right away", "right here", "right the ship"
	case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
		resolve = TagADV // "turned right", "veered right"
	case t.HasTag(TagNOUN) && resolvedAs(prev, TagDET) && !next.HasTag(TagNOUN|TagADJ|TagPROPN):
		resolve = TagNOUN // "the right (to vote)"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+right"
	}
}
