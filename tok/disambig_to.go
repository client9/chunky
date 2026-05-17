package tok

import (
	"strings"
)

// toInfinitiveVerbs is the set of verbs (lowercased) that take a to-infinitive
// complement in their most common usage. Motion and PP-taking verbs like "went",
// "came", "moved" are excluded — they take prepositional "to", not infinitives.
var toInfinitiveVerbs = map[string]bool{
	// refusal / agreement
	"declined": true, "refused": true, "agreed": true, "offered": true,
	// decision / intention
	"decided": true, "chose": true, "intended": true, "aimed": true,
	"planned": true, "sought": true, "attempted": true, "tried": true,
	// ability / effort
	"managed": true, "failed": true, "struggled": true, "continued": true,
	"began": true, "started": true, "stopped": true,
	// desire / expectation
	"wanted": true, "needed": true, "hoped": true, "expected": true,
	"asked": true, "preferred": true,
}

// isToInfinitiveVerb reports whether word is a verb that takes a to-infinitive
// complement ("declined to comment") rather than a prepositional "to" ("went to war").
func isToInfinitiveVerb(word string) bool {
	return toInfinitiveVerbs[strings.ToLower(word)]
}

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
	for i := range tokens {
		disambiguateTo(tokens, i)
	}
	return tokens
}

func disambiguateTo(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADP) || !t.HasTag(TagPART) {
		return
	}
	prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagVERB) && !next.HasTag(TagNOUN):
		resolve = TagPART // "want to run", "decided to leave"
	case next.HasTag(TagVERB) && next.HasTag(TagNOUN) && isToInfinitiveVerb(prev.Word):
		// "declined to comment", "refused to answer": PART, and the following
		// NOUN|VERB is the infinitive — resolve it immediately so the
		// NOUN|VERB-before-PUNCT context rule can't fire first.
		tokens[i+1].Tags = TagVERB
		tokens[i+1].Rule = next.Rule + "+after-to"
		resolve = TagPART
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
