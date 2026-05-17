package tok

// DisambiguateThere resolves the PRON/ADV ambiguity on "there"/"There" tokens.
//
// Existential "there" (followed immediately by AUX or VERB: "there is", "there are")
// is tagged PRON. Locative "there" ("go there", "over there") is tagged ADV.
func DisambiguateThere(tokens []Token) []Token {
	for i := range tokens {
		disambiguateThere(tokens, i)
	}
	return tokens
}

func disambiguateThere(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagPRON) || !t.HasTag(TagADV) {
		return
	}
	tag := TagADV
	if tokenAt(tokens, i+1).HasTag(TagAUX | TagVERB) {
		tag = TagPRON
	}
	tokens[i].Tags = tag
	tokens[i].Rule = t.Rule + "+there"
}
