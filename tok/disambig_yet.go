package tok

// DisambiguateYet resolves the ADV/CCONJ ambiguity on "yet".
//
// ADV (overwhelmingly dominant):
//   - next=ADV|ADJ|DET|NUM → ADV  ("yet another", "not yet complete")
//   - prev=AUX|VERB        → ADV  ("has not yet", "arrived yet")
//
// CCONJ left unresolved: clause-level adversative ("weak, yet compelling")
// requires wider context than single-token lookahead.
func DisambiguateYet(tokens []Token) []Token {
	for i := range tokens {
		disambiguateYet(tokens, i)
	}
	return tokens
}

func disambiguateYet(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADV) || !t.HasTag(TagCCONJ) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagADV | TagADJ | TagDET | TagNUM | TagPART):
		resolve = TagADV
	case prev.HasTag(TagAUX|TagVERB) || resolvedAs(prev, TagPART):
		resolve = TagADV
	case next.HasTag(TagPUNCT):
		resolve = TagADV // "not yet.", "not decided yet."
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+yet"
	}
}
