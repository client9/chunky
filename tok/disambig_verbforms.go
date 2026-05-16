package tok

import "strings"

// DisambiguateVerbForms resolves common NOUN|VERB words that appear frequently
// as finite verbs in subject-verb position. When the preceding token is a
// resolved nominal (NOUN, PROPN, PRON, NUM) the word is a finite verb.
// Also handles quotative inversion: PUNCT says PROPN/NOUN → verb.
func DisambiguateVerbForms(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagVERB) || !t.HasTag(TagNOUN) {
			continue
		}
		lw := strings.ToLower(t.Word)
		switch lw {
		case "says", "say", "remains", "calls", "rose", "fell", "runs", "turns",
			"holds", "needs", "wants", "plans", "shows", "leads", "leaves",
			"means", "takes", "makes", "comes", "goes", "gives", "brings",
			"adds", "argues":
		default:
			continue
		}
		prev := tokenAt(tokens, i-1)
		if prev.HasTag(TagNOUN) || prev.HasTag(TagPROPN) || prev.HasTag(TagPRON) || prev.HasTag(TagNUM) {
			tokens[i].Tags = TagVERB
			tokens[i].Rule = t.Rule + "+verbform"
			continue
		}
		// Quotative inversion: "...", says Bonita → PUNCT before, nominal after.
		if prev.HasTag(TagPUNCT) {
			next := tokenAt(tokens, i+1)
			if next.HasTag(TagNOUN) || next.HasTag(TagPROPN) || next.HasTag(TagPRON) {
				tokens[i].Tags = TagVERB
				tokens[i].Rule = t.Rule + "+verbform"
			}
		}
	}
	return tokens
}
