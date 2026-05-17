package tok

// DisambiguateFree resolves "free" ({ADJ,ADV,VERB}).
// ADJ dominates in almost every context (97–99%).
// VERB only when followed by a direct object.
//
//   - next=PRON|DET → VERB  ("free them", "free the prisoners")
//   - next=NOUN|PROPN|PUNCT|ADP|CCONJ|PART|ADJ|AUX|ADV|SCONJ → ADJ
func DisambiguateFree(tokens []Token) []Token {
	for i := range tokens {
		disambiguateFree(tokens, i)
	}
	return tokens
}

func disambiguateFree(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagVERB) {
		return
	}
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagPRON | TagDET):
		resolve = TagVERB
	case next.HasTag(TagNOUN | TagPROPN | TagPUNCT | TagADP | TagCCONJ | TagPART | TagADJ | TagAUX | TagADV | TagSCONJ):
		resolve = TagADJ
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+free"
	}
}
