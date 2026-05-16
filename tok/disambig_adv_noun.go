package tok

import "strings"

// DisambiguateAdvNoun resolves frequent {ADV,NOUN} and {ADJ,ADV,DET} words.
//
// {ADV,NOUN}:
//   way:   next=ADV|ADJ → ADV ("way too much", "way ahead"); else → NOUN
//   brand: next=ADJ     → ADV ("brand new", "brand fresh"); else → NOUN
//   lot:                → NOUN (always; "a lot of", "lots of")
//
// {ADJ,ADV,DET}:
//   only:  next=VERB|AUX|ADV → ADV; prev=DET → ADJ ("the only one"); next=ADJ → ADV
//   little: next=NOUN → DET; next=ADJ|ADV → ADV
func DisambiguateAdvNoun(tokens []Token) []Token {
	for i, t := range tokens {
		lw := strings.ToLower(t.Word)
		switch lw {
		case "way", "brand", "lot", "only", "little":
		default:
			continue
		}
		prev := tokenAt(tokens, i-1)
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch lw {
		case "way":
			if !t.HasTag(TagADV) || !t.HasTag(TagNOUN) {
				continue
			}
			switch {
			case resolvedAs(prev, TagDET):
				// "the way", "a way" — always nominal
				resolve = TagNOUN
			case next.HasTag(TagADV | TagADJ):
				// "way too much", "way ahead" — intensifier ADV
				resolve = TagADV
			default:
				resolve = TagNOUN
			}
		case "brand":
			if !t.HasTag(TagADV) || !t.HasTag(TagNOUN) {
				continue
			}
			if next.HasTag(TagADJ) {
				resolve = TagADV
			} else {
				resolve = TagNOUN
			}
		case "lot":
			if !t.HasTag(TagADV) || !t.HasTag(TagNOUN) {
				continue
			}
			resolve = TagNOUN
		case "only":
			if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
				continue
			}
			switch {
			case resolvedAs(prev, TagDET) && next.HasTag(TagNOUN|TagPROPN|TagPRON|TagNUM):
				// prenominal ADJ takes priority: "the only person", "the only one"
				resolve = TagADJ
			case next.HasTag(TagVERB | TagAUX | TagADV | TagADP | TagSCONJ | TagPART | TagPRON | TagDET | TagNUM | TagPUNCT):
				// ADV modifying clause-level elements
				resolve = TagADV
			case next.HasTag(TagADJ) && !next.HasTag(TagNOUN):
				resolve = TagADV
			}
		case "little":
			if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
				continue
			}
			switch {
			case next.HasTag(TagADJ | TagADV):
				resolve = TagADV
			case next.HasTag(TagNOUN | TagPROPN):
				resolve = TagDET
			}
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+adv-noun"
		}
	}
	return tokens
}
