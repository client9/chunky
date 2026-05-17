package tok

import "strings"

// possessivePronouns is the closed set of English genitive/possessive pronouns.
// These are all tagged PRON, identical to subject pronouns, so we must
// match by word to distinguish "his own X" (ADJ) from "they own X" (VERB).
var possessivePronouns = map[string]bool{
	"my": true, "your": true, "his": true, "her": true,
	"its": true, "our": true, "their": true,
}

// DisambiguateOwn resolves {ADJ,VERB} words that are almost always ADJ
// when preceded by a possessive pronoun or DET and followed by a noun phrase.
//
// own/live/separate/complete/correct/dry/warm/smooth/secure/frequent/lasting:
//   - prev=possessive-PRON → ADJ  ("his own X", "their live X")
//   - next=AUX → VERB  ("they own/live/complete ...")
func DisambiguateOwn(tokens []Token) []Token {
	for i := range tokens {
		disambiguateOwn(tokens, i)
	}
	return tokens
}

func disambiguateOwn(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagVERB) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case possessivePronouns[strings.ToLower(prev.Word)]:
		resolve = TagADJ // "his own X", "their live broadcast"
	case next.HasTag(TagNOUN|TagPROPN) && !next.HasTag(TagVERB|TagAUX) &&
		prev.HasTag(TagDET|TagADJ|TagNUM):
		resolve = TagADJ // "a dry run", "the complete list", "medium sized dog"
	case next.HasTag(TagAUX):
		resolve = TagVERB // "they separate ... AUX"
	case resolvedAs(prev, TagADV) && next.HasTag(TagNOUN|TagPROPN):
		resolve = TagADJ // "completely separate issue"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+own"
	}
}
