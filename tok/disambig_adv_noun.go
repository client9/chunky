package tok

import "strings"

// DisambiguateAdvNoun resolves frequent {ADV,NOUN} words.
//
// way:   next=ADV|ADJ → ADV ("way too much", "way ahead"); else → NOUN
// brand: next=ADJ     → ADV ("brand new", "brand fresh"); else → NOUN
// lot:                → NOUN (always; "a lot of", "lots of")
func DisambiguateAdvNoun(tokens []Token) []Token {
	for i := range tokens {
		disambiguateAdvNoun(tokens, i)
	}
	return tokens
}

func disambiguateAdvNoun(tokens []Token, i int) {
	t := tokens[i]
	lw := strings.ToLower(t.Word)
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch lw {
	case "way":
		if !t.HasTag(TagADV) || !t.HasTag(TagNOUN) {
			return
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
			return
		}
		if next.HasTag(TagADJ) {
			resolve = TagADV
		} else {
			resolve = TagNOUN
		}
	case "lot":
		if !t.HasTag(TagADV) || !t.HasTag(TagNOUN) {
			return
		}
		resolve = TagNOUN
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adv-noun"
	}
}
