package tok

// DisambiguateBack resolves "back" ({ADV,NOUN}).
//
//   - prev=VERB            → ADV  ("went back", "came back", "step back")
//   - next=NOUN, prev≠VERB → ADJ  ("back door", "back seat")
func DisambiguateBack(tokens []Token) []Token {
	for i := range tokens {
		disambiguateBack(tokens, i)
	}
	return tokens
}

func disambiguateBack(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADV) || !t.HasTag(TagNOUN) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	switch {
	case prev.HasTag(TagVERB):
		tokens[i].Tags = TagADV
		tokens[i].Rule = t.Rule + "+back"
	case next.HasTag(TagNOUN|TagADJ|TagPROPN) && !prev.HasTag(TagVERB):
		tokens[i].Tags = TagADJ
		tokens[i].Rule = t.Rule + "+back"
	}
}
