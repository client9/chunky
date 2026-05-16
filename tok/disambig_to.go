package tok

import "strings"

// DisambiguateTo resolves the ADP/PART ambiguity on "to" and "To".
//
// PART (infinitive marker): "want to run", "decided to leave"
//   - next=VERB (not also NOUN) → PART
//
// ADP (preposition): "go to the store", "according to him"
//   - next=DET|PRON|PROPN|NUM  → ADP (noun phrase follows)
//   - resolvedAs(next, NOUN)    → ADP (pure noun follows)
//
// Note: remaining ambiguous cases (e.g. next={NOUN,VERB}) are resolved
// post-chunk by DisambiguateByChunk (VP → PART, PP → ADP).
func DisambiguateTo(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagADP) || !t.HasTag(TagPART) {
			continue
		}
		if strings.ToLower(t.Word) != "to" {
			continue
		}
		prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
		var resolve Tag
		switch {
		case next.HasTag(TagVERB) && !next.HasTag(TagNOUN):
			resolve = TagPART // "want to run", "decided to leave"
		case next.HasTag(TagDET | TagPRON | TagPROPN | TagNUM):
			resolve = TagADP // "to the store", "to him", "to London", "to five"
		case resolvedAs(next, TagNOUN):
			resolve = TagADP // "to war", "to victory" (pure noun)
		case resolvedAs(prev, TagAUX):
			resolve = TagPART // "should to", "need to", "have to", "going to"
		case next.HasTag(TagADP|TagADV) && !next.HasTag(TagNOUN|TagVERB):
			resolve = TagADP // "to out", "to far" — prepositional use
		case next.HasTag(TagPUNCT):
			resolve = TagADP // sentence-ending "to." is prepositional
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+to"
		}
	}
	return tokens
}
