package tok

// DisambiguateSuch resolves the ADJ/ADV/PART ambiguity on "such" and "Such".
//
// ADJ (prenominal): "such a case", "such cases", "such behavior"
// "such as" is already merged to ADP by MergeLexical — those tokens never reach here.
func DisambiguateSuch(tokens []Token) []Token {
	for i := range tokens {
		disambiguateSuch(tokens, i)
	}
	return tokens
}

func disambiguateSuch(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) {
		return
	}
	next := tokenAt(tokens, i+1)
	if next.HasTag(TagDET | TagNOUN | TagPROPN | TagADJ) {
		tokens[i].Tags = TagADJ
		tokens[i].Rule = t.Rule + "+such"
	}
}
