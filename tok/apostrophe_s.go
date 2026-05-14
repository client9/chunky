package tok

import (
	"strings"

	"github.com/client9/chunky"
)

// auxHosts is the closed set of words whose "'s" contraction is a copula/auxiliary,
// not a possessive. All others (nouns, proper nouns, indefinite pronouns) are PART.
var auxHosts = map[string]bool{
	"he": true, "she": true, "it": true,
	"that": true, "this": true,
	"what": true, "where": true, "who": true, "how": true,
	"there": true, "here": true,
}

// DisambiguateApostropheS resolves the AUX/PART ambiguity on "'s" tokens by
// inspecting the surface form of the immediately preceding token.
func DisambiguateApostropheS(tokens []Token) []Token {
	for i, t := range tokens {
		if t.Word != "'s" {
			continue
		}
		tag := chunky.Tag(chunky.TagPART)
		if i > 0 && auxHosts[strings.ToLower(tokens[i-1].Word)] {
			tag = chunky.TagAUX
		}
		tokens[i].Tags = []chunky.Tag{tag}
		tokens[i].Rule = "apostrophe-s"
	}
	return tokens
}
