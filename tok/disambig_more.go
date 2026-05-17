package tok


// DisambiguateMore resolves the ADV/DET ambiguity on degree quantifiers.
//
// Quantifier/determiner (DET): "more people", "most cities", "much time"
// Degree modifier (ADV):       "more quickly", "most important", "much faster"
//
// Resolved cases:
//   - next=ADJ|ADV → ADV  (degree modifier before comparative/superlative)
//   - next=NOUN    → DET  (quantifier before noun)
func DisambiguateMore(tokens []Token) []Token {
	for i := range tokens {
		disambiguateMore(tokens, i)
	}
	return tokens
}

func disambiguateMore(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADV) || !t.HasTag(TagDET) {
		return
	}
	next := tokenAt(tokens, i+1)
	prev := tokenAt(tokens, i-1)
	var resolve Tag
	switch {
	case next.HasTag(TagADJ | TagADV):
		resolve = TagADV
	case next.HasTag(TagNOUN) && !next.HasTag(TagVERB):
		resolve = TagDET
	case next.Word == "of":
		resolve = TagDET // "most of the time", "more of the same"
	case next.HasTag(TagPUNCT | TagCCONJ):
		resolve = TagADV // "they want more.", "more and more"
	case prev.HasTag(TagCCONJ):
		resolve = TagADV // "or more", "or less"
	case prev.HasTag(TagAUX|TagVERB) && !next.HasTag(TagNOUN|TagADJ|TagPROPN):
		resolve = TagADV // "is more", "becomes more" (not before noun phrase)
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+more"
	}
}
