package tok

import "strings"

// DisambiguateQuantifiers resolves collective/distributive determiners.
//
// both/neither: next=DET|NOUN|ADJ|PROPN|NUM → DET; prev=PRON|VERB → PRON
// either:       next=DET|NOUN|ADJ|PROPN     → DET
// all:          next=DET|NOUN|ADJ|PROPN|NUM → DET; prev=PRON|VERB → PRON
// each:         next=NOUN|PROPN|ADJ|NUM     → DET; next=ADP|AUX|VERB|PUNCT|CCONJ → PRON
// any:          next=NOUN|PROPN|ADJ|NUM     → DET; next=ADP|PUNCT → PRON
func DisambiguateQuantifiers(tokens []Token) []Token {
	for i := range tokens {
		disambiguateQuantifiers(tokens, i)
	}
	return tokens
}

func disambiguateQuantifiers(tokens []Token, i int) {
	t := tokens[i]
	lw := strings.ToLower(t.Word)
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch lw {
	case "both", "neither":
		if !t.HasTag(TagDET) || !t.HasTag(TagCCONJ) {
			return
		}
		switch {
		case next.HasTag(TagNOUN|TagPROPN|TagDET|TagADJ|TagNUM) && !next.HasTag(TagVERB|TagAUX):
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
			return
		}
		switch {
		case next.HasTag(TagDET|TagNOUN|TagADJ|TagPROPN) && !next.HasTag(TagVERB|TagAUX):
			resolve = TagDET
		case resolvedAs(prev, TagPART) || resolvedAs(prev, TagADV):
			resolve = TagADV // "not either", "nor either"
		case next.HasTag(TagADP):
			resolve = TagADV // "either of them"
		}
	case "all":
		if !t.HasTag(TagDET) {
			return
		}
		switch {
		case next.HasTag(TagDET|TagNOUN|TagADJ|TagPROPN|TagNUM) && !next.HasTag(TagVERB|TagAUX):
			resolve = TagDET
		case resolvedAs(prev, TagPRON) || resolvedAs(prev, TagNOUN):
			resolve = TagPRON
		case prev.HasTag(TagADP|TagADV) && next.HasTag(TagPUNCT):
			resolve = TagADV // "above all.", "after all,", "not at all."
		}
	case "each":
		if !t.HasTag(TagDET) {
			return
		}
		switch {
		case next.HasTag(TagNOUN|TagPROPN|TagADJ|TagNUM) && !next.HasTag(TagVERB|TagAUX):
			resolve = TagDET // "each team", "each player"
		case next.HasTag(TagADP | TagAUX | TagVERB | TagPUNCT | TagCCONJ):
			resolve = TagPRON // "each of", "each will", "each said", "each.", "each and"
		}
	case "any":
		if !t.HasTag(TagDET) {
			return
		}
		switch {
		case next.HasTag(TagNOUN|TagPROPN|TagADJ|TagNUM) && !next.HasTag(TagVERB|TagAUX):
			resolve = TagDET // "any team", "any suggestion"
		case next.HasTag(TagADP):
			resolve = TagPRON // "any of them"
		case next.HasTag(TagPUNCT):
			resolve = TagPRON // "if any.", "not any,"
		}
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+quant"
	}
}
