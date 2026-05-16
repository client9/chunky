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
//
// It also resolves NOUN/VERB ambiguity on the neighbors of a possessive "'s":
//   - The possessor (token before "'s") is always a noun: "father 's"
//   - The head of the possessive phrase (token after "'s") is always a noun: "'s board"
//
// These two neighbor resolutions are done here rather than in the context rule
// table because PART conflates possessive "'s" and infinitival "to". Their
// corpus statistics cancel out, so no corpus-derived rule clears the 10× ratio
// threshold. The possessive case is a linguistic axiom, not a probabilistic rule.
func DisambiguateApostropheS(tokens []Token) []Token {
	for i, t := range tokens {
		if t.Word != "'s" {
			continue
		}
		tag := chunky.Tag(chunky.TagPART)
		if i > 0 && auxHosts[strings.ToLower(tokens[i-1].Word)] {
			tag = chunky.TagAUX
		}
		tokens[i].Tags = tag
		tokens[i].Rule = "apostrophe-s"

		if tag == chunky.TagPART {
			// Possessor: the token before a possessive "'s" is always a noun.
			if i > 0 && tokens[i-1].HasTag(chunky.TagNOUN) && tokens[i-1].HasTag(chunky.TagVERB) {
				tokens[i-1].Tags = chunky.TagNOUN
				tokens[i-1].Rule = tokens[i-1].Rule + "+poss-host"
			}
			// Possessed head: the token after a possessive "'s" is always a noun.
			if i+1 < len(tokens) && tokens[i+1].HasTag(chunky.TagNOUN) && tokens[i+1].HasTag(chunky.TagVERB) {
				tokens[i+1].Tags = chunky.TagNOUN
				tokens[i+1].Rule = tokens[i+1].Rule + "+poss-head"
			}
		}
	}
	return tokens
}
