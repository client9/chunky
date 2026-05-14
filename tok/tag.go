package tok

import (
	"strings"

	"github.com/client9/chunky"
)

// LexicalTag assigns candidate tags to untagged tokens by looking up each
// token's lowercase form in the compiled lexicon and the runtime AbbreviationTags
// map. Tokens that already carry candidates (e.g. compound tokens from
// MergeLexical) are left unchanged.
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
		var candidates []chunky.Tag
		if len(src) > 0 {
			candidates = make([]chunky.Tag, len(src))
			copy(candidates, src)
		}
		tokens[i].Candidates = candidates
		tokens[i].Rule = rule
	}
	return tokens
}
