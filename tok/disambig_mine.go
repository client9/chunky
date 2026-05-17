package tok

// DisambiguateMine resolves the NOUN/PRON ambiguity on "mine".
//
// mine — PRON (possessive) or NOUN (excavation):
//   - prev=AUX|VERB → PRON  ("is mine", "was mine")
//   - prev=DET       → NOUN  ("the mine", "a mine")
//   - next=NOUN|PROPN → NOUN  ("mine shaft")
//   - next=ADP       → NOUN  ("mine of gold")
func DisambiguateMine(tokens []Token) []Token {
	for i := range tokens {
		disambiguateMine(tokens, i)
	}
	return tokens
}

func disambiguateMine(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagNOUN) || !t.HasTag(TagPRON) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case prev.HasTag(TagAUX | TagVERB):
		resolve = TagPRON
	case resolvedAs(prev, TagDET):
		resolve = TagNOUN
	case next.HasTag(TagNOUN | TagPROPN):
		resolve = TagNOUN
	case next.HasTag(TagADP):
		resolve = TagNOUN
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+mine"
	}
}
