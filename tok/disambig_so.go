package tok


// DisambiguateSo resolves the ADV/SCONJ ambiguity on "so".
//
// Intensifier (ADV):    "so good", "so quickly"
// Subordinator (SCONJ): "so that everyone could see"
//
// Resolved cases:
//   - next=ADJ|ADV → ADV   ("so good", "so carefully")
//   - next=SCONJ   → SCONJ ("so that")
//
// Note: "so that" is handled upstream by MergeLexical as a compound SCONJ token;
// the next=SCONJ rule catches remaining cases like "so as", "so because".
func DisambiguateSo(tokens []Token) []Token {
	for i := range tokens {
		disambiguateSo(tokens, i)
	}
	return tokens
}

func disambiguateSo(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADV) || !t.HasTag(TagSCONJ) {
		return
	}
	prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagADJ | TagADV):
		resolve = TagADV
	case resolvedAs(next, TagSCONJ):
		resolve = TagSCONJ
	case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
		resolve = TagADV // "did so", "think so"
	case prev.HasTag(TagCCONJ):
		resolve = TagADV // "and so on", "or so", "but so"
	case next.HasTag(TagADP):
		resolve = TagADV // "so on", "so far", "so to speak"
	case next.HasTag(TagPRON | TagDET):
		resolve = TagSCONJ // "so he left", "so the war ended"
	case next.HasTag(TagAUX):
		resolve = TagSCONJ // "so was", "so had", "so does" (inverted/fronted)
	case resolvedAs(next, TagNOUN) || resolvedAs(next, TagPROPN):
		resolve = TagSCONJ // "so peace prevailed", "so victory was"
	case next.HasTag(TagPUNCT):
		resolve = TagADV // sentence-final "so."
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+so"
	}
}
