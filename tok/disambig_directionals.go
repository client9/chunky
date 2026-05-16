package tok

import "strings"

// DisambiguateDirectionals resolves the ADJ/ADV/NOUN ambiguity on cardinal direction words.
//
// "south" lexicon: ADJ|ADV|NOUN
//   - next=ADP → ADV  ("south of the border", "heading south of")
//   - prev=DET → NOUN ("the south", "a south")
//
// "north", "east", "west" lexicon: ADJ|NOUN
//   - prev=DET → NOUN ("the north", "in the east")
//   otherwise → ADJ  ("north side", "east coast")
func DisambiguateDirectionals(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagADJ) || !t.HasTag(TagNOUN) {
			continue
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
		case "north", "east", "west":
			switch {
			case next.HasTag(TagPROPN):
				resolve = TagADJ // "North Korea", "East Germany", "West Africa"
			case detPrevNoun:
				resolve = TagNOUN
			}
		default:
			continue
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+dir"
		}
	}
	return tokens
}
