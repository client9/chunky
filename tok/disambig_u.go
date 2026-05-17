package tok

// DisambiguateU resolves "U"/"u" ({NOUN,PRON}).
// The Brown corpus over-tags this letter/abbreviation as PRON; it is almost
// always NOUN (letter, shape, abbreviation: "U-turn", "vitamin U").
//
//   - prev=DET|ADJ|NUM → NOUN  ("a U shape", "3 U bolts")
//   - next=NOUN|PROPN   → NOUN  ("U shape", "U bolt")
//   - next=PUNCT         → NOUN  (letter reference: "vitamin U.")
func DisambiguateU(tokens []Token) []Token {
	for i := range tokens {
		disambiguateU(tokens, i)
	}
	return tokens
}

func disambiguateU(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagNOUN) || !t.HasTag(TagPRON) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case prev.HasTag(TagDET | TagADJ | TagNUM):
		resolve = TagNOUN
	case next.HasTag(TagNOUN | TagPROPN):
		resolve = TagNOUN
	case next.HasTag(TagPUNCT):
		resolve = TagNOUN
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+mine"
	}
}
