package tok

import "strings"

// DisambiguateFollowing resolves the ADP/NOUN/VERB ambiguity on "following".
//
// VERB (gerund/participle, used prepositionally):
//   - next=DET|PROPN|PRON|NUM|ADJ → VERB  ("following the announcement", "following his resignation")
//
// NOUN ("the following"):
//   - next=AUX → NOUN  ("the following was announced")
//
// ADJ ("the following day") — left ambiguous with VERB when next=NOUN (50/50 in corpus).
func DisambiguateFollowing(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagVERB) || !t.HasTag(TagNOUN) {
			continue
		}
		lw := strings.ToLower(t.Word)
		if lw != "following" {
			continue
		}
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch {
		case next.HasTag(TagDET | TagPROPN | TagPRON | TagNUM | TagADJ):
			resolve = TagVERB
		case next.HasTag(TagAUX):
			resolve = TagNOUN
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+following"
		}
	}
	return tokens
}
