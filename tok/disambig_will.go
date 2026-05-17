package tok

// DisambiguateWill resolves the AUX/NOUN ambiguity on "will" and "Will".
//
// Modal use (AUX): "will go", "will not", "will never", "will he?"
// Noun use (NOUN): "his will", "the will of the people"
func DisambiguateWill(tokens []Token) []Token {
	for i := range tokens {
		disambiguateWill(tokens, i)
	}
	return tokens
}

func disambiguateWill(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagAUX) || !t.HasTag(TagNOUN) {
		return
	}
	prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case resolvedAs(prev, TagDET):
		resolve = TagNOUN
	case next.HasTag(TagVERB | TagAUX | TagADV | TagPART | TagPRON):
		resolve = TagAUX
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+will"
	}
}
