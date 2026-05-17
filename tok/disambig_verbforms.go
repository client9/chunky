package tok

import "strings"

// possessivePRON returns true for pronouns that introduce NPs rather than
// serving as VP subjects (my, your, his, her, its, our, their, whose).
func possessivePRON(t Token) bool {
	if !resolvedAs(t, TagPRON) {
		return false
	}
	switch strings.ToLower(t.Word) {
	case "my", "your", "his", "her", "its", "our", "their", "whose":
		return true
	}
	return false
}

// DisambiguateVerbForms resolves common NOUN|VERB words that appear frequently
// as finite verbs in subject-verb position. When the preceding token is a
// resolved nominal (NOUN, PROPN, PRON, NUM) the word is a finite verb.
// Also handles quotative inversion: PUNCT says PROPN/NOUN → verb.
func DisambiguateVerbForms(tokens []Token) []Token {
	for i := range tokens {
		disambiguateVerbForms(tokens, i)
	}
	return tokens
}

func disambiguateVerbForms(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagVERB) || !t.HasTag(TagNOUN) {
		return
	}
	prev := tokenAt(tokens, i-1)
	// Possessive pronouns (his plans, her calls) introduce NPs — don't treat as subject.
	if possessivePRON(prev) {
		return
	}
	if prev.HasTag(TagNOUN) || prev.HasTag(TagPROPN) || prev.HasTag(TagPRON) || prev.HasTag(TagNUM) {
		tokens[i].Tags = TagVERB
		tokens[i].Rule = t.Rule + "+verbform"
		return
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
