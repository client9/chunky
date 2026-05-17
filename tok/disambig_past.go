package tok

// DisambiguatePast resolves the ADJ/ADP/NOUN ambiguity on "past".
//
// "past" (ADJ 45%, ADP 30%, NOUN 25%):
//   - next=NOUN|NUM|ADJ → ADJ  ("past performance", "past 10 years", "past president")
//   - next=PUNCT|AUX    → NOUN ("in the past.", "the past was")
//   - next=DET|PROPN    → ADP  ("past the building", "past Paris")
func DisambiguatePast(tokens []Token) []Token {
	for i := range tokens {
		disambiguatePast(tokens, i)
	}
	return tokens
}

func disambiguatePast(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) {
		return
	}
	if !t.HasTag(TagADP) && !t.HasTag(TagNOUN) {
		return
	}
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagNOUN | TagNUM | TagADJ):
		resolve = TagADJ
	case next.HasTag(TagPUNCT | TagAUX):
		resolve = TagNOUN
	case next.HasTag(TagDET | TagPROPN):
		resolve = TagADP
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+past"
	}
}
