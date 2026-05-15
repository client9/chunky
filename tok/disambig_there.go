package tok

import (
	"strings"

	"github.com/client9/chunky"
)

// DisambiguateThere resolves the PRON/ADV ambiguity on "there"/"There" tokens.
//
// Existential "there" (followed immediately by AUX or VERB: "there is", "there are")
// is tagged PRON. Locative "there" ("go there", "over there") is tagged ADV.
func DisambiguateThere(tokens []Token) []Token {
	for i, t := range tokens {
		if strings.ToLower(t.Word) != "there" {
			continue
		}
		if !t.HasTag(chunky.TagPRON) || !t.HasTag(chunky.TagADV) {
			continue
		}
		tag := chunky.Tag(chunky.TagADV)
		if i+1 < len(tokens) {
			next := tokens[i+1]
			if next.HasTag(chunky.TagAUX) || next.HasTag(chunky.TagVERB) {
				tag = chunky.TagPRON
			}
		}
		tokens[i].Tags = []chunky.Tag{tag}
		tokens[i].Rule = t.Rule + "+there"
	}
	return tokens
}
