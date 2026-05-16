package tok

import (
	"strings"

	"github.com/client9/chunky"
)

// LexicalTag assigns tags to untagged tokens by looking up each token's
// lowercase form in the compiled lexicon and the runtime AbbreviationTags map.
// Tokens that already carry tags (e.g. compound tokens from chunky.MergeLexical)
// are left unchanged.
func LexicalTag(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.IsUnknownTag() {
			continue
		}
		lower := strings.ToLower(t.Word)
		rule := "lexicon"
		var tags chunky.Tag
		// WordTags is hand-curated and overrides the compiled lexicon.
		if t, ok := chunky.WordTags[lower]; ok {
			tags = t
			rule = "words"
		} else if t := wordtagmap[lower]; t != 0 {
			tags = t
		} else if t, ok := chunky.ClosedFormTags[lower]; ok {
			tags = t
			rule = "closed"
		} else if t, ok := chunky.AbbreviationTags[lower]; ok {
			tags = t
			rule = "abbrev"
		} else {
			rule = ""
		}
		tokens[i].Tags = tags
		tokens[i].Rule = rule
	}
	return tokens
}
