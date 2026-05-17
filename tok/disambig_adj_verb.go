package tok

// DisambiguateAdjVerb resolves the generic ADJ/VERB ambiguity using structural
// context. This covers past participles and dual-category words not handled by
// more specific disambiguators (experienced, sized, thin, wet, learned, etc.).
//
// ADJ (prenominal or predicative):
//   - prev=DET|ADJ|NUM and next=NOUN|PROPN (unambiguous) → ADJ
//   - possessive pronoun + word → ADJ  (handled by DisambiguateOwn)
//
// VERB (finite or participial in VP):
//   - prev=AUX → VERB  ("was experienced", "is marked", "has varied")
//   - prev=PRON|NOUN|PROPN and next=DET|NOUN|PRON → VERB  (SVO: "they lay plans")
func DisambiguateAdjVerb(tokens []Token) []Token {
	for i := range tokens {
		disambiguateAdjVerb(tokens, i)
	}
	return tokens
}

func disambiguateAdjVerb(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagVERB) {
		return
	}
	if t.IsResolved() {
		return
	}
	prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case prev.HasTag(TagAUX):
		resolve = TagVERB // "was experienced", "is marked", "has varied"
	case prev.HasTag(TagDET|TagADJ|TagNUM|TagADV) && next.HasTag(TagNOUN|TagPROPN) && !next.HasTag(TagVERB|TagAUX):
		resolve = TagADJ // "the experienced pilot", "highly experienced team", "a dry run"
	case resolvedAs(prev, TagDET) && !next.HasTag(TagNOUN|TagADJ|TagPROPN):
		resolve = TagADJ // "the complete.", "a warm" — standalone after DET
	case prev.HasTag(TagPRON|TagNOUN|TagPROPN) && !prev.HasTag(TagADJ) &&
		next.HasTag(TagDET|TagPRON|TagNUM) && !next.HasTag(TagVERB):
		resolve = TagVERB // "they separate [the groups]", "birds lay [eggs]"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adjverb"
	}
}
