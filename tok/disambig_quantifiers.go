package tok

import "strings"

// DisambiguateQuantifiers resolves collective/distributive determiners.
//
// both/neither: next=DET|NOUN|ADJ|PROPN|NUM → DET; prev=PRON|VERB → PRON
// either:       next=DET|NOUN|ADJ|PROPN     → DET
// all:          next=DET|NOUN|ADJ|PROPN|NUM → DET; prev=PRON|VERB → PRON
func DisambiguateQuantifiers(tokens []Token) []Token {
	for i, t := range tokens {
		lw := strings.ToLower(t.Word)
		switch lw {
		case "both", "neither", "either", "all":
		default:
			continue
		}
		prev := tokenAt(tokens, i-1)
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch lw {
		case "both", "neither":
			if !t.HasTag(TagDET) || !t.HasTag(TagCCONJ) {
				continue
			}
			switch {
			case next.HasTag(TagNOUN | TagPROPN | TagDET | TagADJ | TagNUM):
				// "both teams", "both the teams", "both Germany", "both important"
				resolve = TagDET
			case next.HasTag(TagADP | TagVERB | TagAUX | TagPUNCT | TagADV | TagSCONJ):
				// "both of them", "both were there", "both.", "both when"
				resolve = TagPRON
			case resolvedAs(prev, TagPRON) || prev.HasTag(TagVERB):
				resolve = TagPRON // floating: "they were both present"
			case next.HasTag(TagPRON):
				resolve = TagPRON // "both they" — unusual but PRON context
			}
		case "either":
			if !t.HasTag(TagDET) {
				continue
			}
			switch {
			case next.HasTag(TagDET | TagNOUN | TagADJ | TagPROPN):
				resolve = TagDET
			case resolvedAs(prev, TagPART) || resolvedAs(prev, TagADV):
				resolve = TagADV // "not either", "nor either"
			case next.HasTag(TagADP):
				resolve = TagADV // "either of them"
			}
		case "all":
			if !t.HasTag(TagDET) {
				continue
			}
			switch {
			case next.HasTag(TagDET | TagNOUN | TagADJ | TagPROPN | TagNUM):
				resolve = TagDET
			case resolvedAs(prev, TagPRON) || resolvedAs(prev, TagNOUN):
				resolve = TagPRON
			}
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+quant"
		}
	}
	return tokens
}
