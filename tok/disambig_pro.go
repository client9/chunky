package tok

// DisambiguatePro resolves "pro" ({ADJ,...}).
// In corpus "pro" is overwhelmingly ADJ when followed by a nominal or adjectival head.
//
//   - next=ADJ|NOUN|PROPN → ADJ  ("pro rata", "pro wrestler", "pro sports")
func DisambiguatePro(tokens []Token) []Token {
	for i := range tokens {
		disambiguatePro(tokens, i)
	}
	return tokens
}

func disambiguatePro(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) {
		return
	}
	next := tokenAt(tokens, i+1)
	if next.HasTag(TagADJ | TagNOUN | TagPROPN) {
		tokens[i].Tags = TagADJ
		tokens[i].Rule = t.Rule + "+past"
	}
}
