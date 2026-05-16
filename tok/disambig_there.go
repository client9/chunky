package tok

import "strings"

// DisambiguateThere resolves the PRON/ADV ambiguity on "there"/"There" tokens.
//
// Existential "there" (followed immediately by AUX or VERB: "there is", "there are")
// is tagged PRON. Locative "there" ("go there", "over there") is tagged ADV.
func DisambiguateThere(tokens []Token) []Token {
	for i, t := range tokens {
		if strings.ToLower(t.Word) != "there" {
			continue
		}
		if !t.HasTag(TagPRON) || !t.HasTag(TagADV) {
			continue
		}
		tag := TagADV
		if tokenAt(tokens, i+1).HasTag(TagAUX | TagVERB) {
			tag = TagPRON
		}
		tokens[i].Tags = tag
		tokens[i].Rule = t.Rule + "+there"
	}
	return tokens
}
