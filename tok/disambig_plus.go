package tok


// DisambiguatePlus resolves {CCONJ,PROPN} for "plus".
//
// CCONJ (additive):  "three plus four", "cost plus tax", "a bonus, plus benefits"
// PROPN:             "Disney Plus", "Apple TV+" — handled by RetagCapitalized
//
//   - next=NUM|NOUN|ADJ|PRON|DET → CCONJ
//   - prev=NUM|NOUN|ADJ|PRON    → CCONJ
func DisambiguatePlus(tokens []Token) []Token {
	for i := range tokens {
		disambiguatePlus(tokens, i)
	}
	return tokens
}

func disambiguatePlus(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagCCONJ) || !t.HasTag(TagPROPN) {
		return
	}
	// Capitalized "Plus" mid-sentence is a brand name (Disney Plus, Apple TV+).
	// Leave it for RetagCapitalized rather than forcing CCONJ.
	if t.Word[0] == 'P' {
		return
	}
	prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagNUM | TagNOUN | TagADJ | TagPRON | TagDET | TagPROPN):
		resolve = TagCCONJ
	case prev.HasTag(TagNUM | TagNOUN | TagADJ | TagPRON | TagPROPN):
		resolve = TagCCONJ
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+plus"
	}
}
