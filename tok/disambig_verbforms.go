package tok


// DisambiguateVerbForms resolves common NOUN|VERB words that appear frequently
// as finite verbs in subject-verb position. When the preceding token is a
// resolved nominal (NOUN, PROPN, PRON, NUM) the word is a finite verb.
// Also handles quotative inversion: PUNCT says PROPN/NOUN → verb.
func DisambiguateVerbForms(tokens []Token) []Token {
	for i := range tokens {
		disambiguateVerbForms(tokens, i)
	}
	return tokens
}

func disambiguateVerbForms(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagVERB) || !t.HasTag(TagNOUN) {
		return
	}
	prev := tokenAt(tokens, i-1)
	if prev.HasTag(TagNOUN) || prev.HasTag(TagPROPN) || prev.HasTag(TagPRON) || prev.HasTag(TagNUM) {
		tokens[i].Tags = TagVERB
		tokens[i].Rule = t.Rule + "+verbform"
		return
	}
	// Quotative inversion: "...", says Bonita → PUNCT before, nominal after.
	if prev.HasTag(TagPUNCT) {
		next := tokenAt(tokens, i+1)
		if next.HasTag(TagNOUN) || next.HasTag(TagPROPN) || next.HasTag(TagPRON) {
			tokens[i].Tags = TagVERB
			tokens[i].Rule = t.Rule + "+verbform"
		}
	}
}
