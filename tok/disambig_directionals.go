package tok

import "strings"

// DisambiguateDirectionals resolves the ADJ/ADV/NOUN ambiguity on cardinal direction words.
//
// "south" lexicon: ADJ|ADV|NOUN
//   - next=ADP → ADV  ("south of the border", "heading south of")
//   - prev=DET → NOUN ("the south", "a south")
//
// "north", "east", "west", "northwest", "northeast", "southeast", "southwest" lexicon: ADJ|NOUN
//   - next=PROPN → ADJ ("North Korea", "Southwest Airlines", "Southeast Asia")
//   - prev=DET   → NOUN ("the north", "the northwest", "in the east")
func DisambiguateDirectionals(tokens []Token) []Token {
	for i := range tokens {
		disambiguateDirectionals(tokens, i)
	}
	return tokens
}

func disambiguateDirectionals(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagNOUN) {
		return
	}
	lw := strings.ToLower(t.Word)
	prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
	var resolve Tag
	// prev=DET only resolves to NOUN when the directional is not prenominal
	// (i.e. not followed by a noun or adjective it would modify).
	detPrevNoun := resolvedAs(prev, TagDET) && !next.HasTag(TagNOUN|TagADJ|TagPROPN)
	switch lw {
	case "south":
		switch {
		case next.HasTag(TagPROPN):
			resolve = TagADJ // "South Korea", "South Africa"
		case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
			resolve = TagADV
		case t.HasTag(TagADV) && resolvedAs(next, TagADP):
			resolve = TagADV
		case detPrevNoun:
			resolve = TagNOUN
		}
	case "north", "east", "west",
		"northwest", "northeast", "southeast", "southwest":
		switch {
		case next.HasTag(TagPROPN):
			resolve = TagADJ // "North Korea", "Northeast Asia", "Southwest Airlines"
		case resolvedAs(next, TagADP):
			resolve = TagNOUN // "northwest of the city", "north of the border"
		case detPrevNoun:
			resolve = TagNOUN // "the northwest", "the southeast"
		}
	default:
		return
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+dir"
	}
}
