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
		src := wordtagmap[lower]
		rule := "lexicon"
		if len(src) == 0 {
			// AbbreviationTags is the runtime-editable override layer; entries
			// added there (e.g. contraction suffixes) don't require regenerating
			// the compiled lexicon.
			if tags, ok := chunky.AbbreviationTags[lower]; ok {
				src = tags
			} else {
				rule = ""
			}
		}
		// Copy to avoid aliasing into the compiled lexicon slice.
		var tags []chunky.Tag
		if len(src) > 0 {
			tags = make([]chunky.Tag, len(src))
			copy(tags, src)
		}
		tokens[i].Tags = tags
		tokens[i].Rule = rule
	}
	return tokens
}
