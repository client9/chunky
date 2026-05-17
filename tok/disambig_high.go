package tok

// DisambiguateHigh resolves "high" ({ADJ,ADV,NOUN}).
//
// ADJ (prenominal or predicative):
//   - next=NOUN|PROPN|ADJ|NUM|CCONJ|PUNCT → ADJ  (94–100%)
//   - prev=DET|ADP|ADV|AUX (resolved)     → ADJ  (96–99%)
//   - prev=AUX + next=PUNCT|CCONJ|SCONJ|ADP → ADJ  ("was high.", "is high,")
//
// ADV (post-verbal or degree modifier):
//   - prev=VERB (not AUX)        → ADV  ("ran high", "went high")
//   - prev=ADV|PART              → ADV  ("very high", "not high")
//   - next=ADV (unambiguous)     → ADV  ("high above", "high enough")
//   - next=ADP (not nominal)     → ADV  ("high on", "high in")
//   - next=PUNCT|CCONJ + prev=ADP|ADV → ADV ("came in high,")
//
// NOUN (only when NOUN bit is set):
//   - prev=NOUN|ADJ + no following nominal → NOUN  ("record high", "all-time high")
func DisambiguateHigh(tokens []Token) []Token {
	for i := range tokens {
		disambiguateHigh(tokens, i)
	}
	return tokens
}

func disambiguateHigh(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	// ADJ cases
	case next.HasTag(TagNOUN | TagPROPN | TagADJ | TagNUM | TagCCONJ | TagPUNCT):
		resolve = TagADJ
	case resolvedAs(prev, TagDET) || resolvedAs(prev, TagADP) || resolvedAs(prev, TagADV) || resolvedAs(prev, TagAUX):
		resolve = TagADJ
	case prev.HasTag(TagAUX) && next.HasTag(TagPUNCT|TagCCONJ|TagSCONJ|TagADP):
		resolve = TagADJ // predicative: "was high.", "is high,"
	// ADV cases (only when ADV bit is set)
	case t.HasTag(TagADV) && prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
		resolve = TagADV // "ran high", "went high"
	case t.HasTag(TagADV) && prev.HasTag(TagADV|TagPART):
		resolve = TagADV // "very high", "not high"
	case t.HasTag(TagADV) && next.HasTag(TagADV) && !next.HasTag(TagNOUN|TagVERB|TagADJ):
		resolve = TagADV // "high above", "high enough"
	case t.HasTag(TagADV) && next.HasTag(TagADP) && !next.HasTag(TagNOUN|TagVERB):
		resolve = TagADV // "high on", "high in"
	case t.HasTag(TagADV) && next.HasTag(TagPUNCT|TagCCONJ) && prev.HasTag(TagADP|TagADV):
		resolve = TagADV // "came in high,", "running high and"
	// NOUN cases (only when NOUN bit is set)
	case t.HasTag(TagNOUN) && prev.HasTag(TagNOUN|TagADJ) && !prev.HasTag(TagVERB) && !next.HasTag(TagNOUN|TagADJ|TagPROPN):
		resolve = TagNOUN // "record high", "all-time high"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+high"
	}
}
