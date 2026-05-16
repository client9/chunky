package tok

import "strings"

// DisambiguateVerbForms resolves common NOUN|VERB words that appear frequently
// as finite verbs in subject-verb position. When the preceding token is a
// resolved nominal (NOUN, PROPN, PRON, NUM) the word is a finite verb.
func DisambiguateVerbForms(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagVERB) || !t.HasTag(TagNOUN) {
			continue
		}
		lw := strings.ToLower(t.Word)
		switch lw {
		case "says", "remains", "calls", "rose", "fell", "runs", "turns",
			"holds", "needs", "wants", "plans", "shows", "leads", "leaves",
			"means", "takes", "makes", "comes", "goes", "gives", "brings":
		default:
			continue
		}
		prev := tokenAt(tokens, i-1)
		if prev.HasTag(TagNOUN) || prev.HasTag(TagPROPN) || prev.HasTag(TagPRON) || prev.HasTag(TagNUM) {
			tokens[i].Tags = TagVERB
			tokens[i].Rule = t.Rule + "+verbform"
		}
	}
	return tokens
}
