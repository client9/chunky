package tok

import "strings"

// DisambiguateOne resolves the DET/NUM/PRON ambiguity on "one".
//
// NUM (dominant): "one day", "one of", "chapter one", "one more"
//   - next=ADP|NOUN|ADJ|NUM|PROPN|DET|CCONJ|PUNCT → NUM  (93–100% in corpus)
//
// PRON ("one" as generic pronoun: "one must consider"):
//   - context-rule pass handles "one of" → PRON already
//   - next=AUX is ambiguous (NUM 17%, PRON 70%) — left unresolved here
func DisambiguateOne(tokens []Token) []Token {
	for i := range tokens {
		disambiguateOne(tokens, i)
	}
	return tokens
}

func disambiguateOne(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagNUM) {
		return
	}
	prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case strings.ToLower(prev.Word) == "no":
		resolve = TagPRON // "no one" → generic pronoun ("nobody")
	case next.HasTag(TagADP | TagNOUN | TagADJ | TagNUM | TagPROPN | TagDET | TagCCONJ | TagPUNCT):
		resolve = TagNUM
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+one"
	}
}
