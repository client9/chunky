package tok

// DisambiguateWell resolves "well" ({ADV,NOUN}).
//
//   - prev=VERB → ADV   ("did well", "worked well")
//   - prev=DET  → NOUN  ("the well", "an oil well")
func DisambiguateWell(tokens []Token) []Token {
	for i := range tokens {
		disambiguateWell(tokens, i)
	}
	return tokens
}

func disambiguateWell(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADV) || !t.HasTag(TagNOUN) {
		return
	}
	prev := tokenAt(tokens, i-1)
	switch {
	case prev.HasTag(TagVERB):
		tokens[i].Tags = TagADV
		tokens[i].Rule = t.Rule + "+well"
	case resolvedAs(prev, TagDET):
		tokens[i].Tags = TagNOUN
		tokens[i].Rule = t.Rule + "+well"
	}
}
